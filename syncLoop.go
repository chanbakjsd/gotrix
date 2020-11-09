package gotrix

import (
	"errors"
	"fmt"
	"time"

	"github.com/chanbakjsd/gotrix/api"
	"github.com/chanbakjsd/gotrix/debug"
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

const (
	eventsToFetchPerRoom = 50
	minBackoffTime       = 1
	maxBackoffTime       = 300
	syncTimeout          = 5000
)

// ErrAlreadyClosed is the error returned by (*Client).Close() when called again.
var ErrAlreadyClosed = errors.New("client already closed")

// DefaultFilter is the default filter used by the client.
var DefaultFilter = event.GlobalFilter{
	Room: event.RoomFilter{
		IncludeLeave: false,
		State: event.StateFilter{
			LazyLoadMembers: true,
		},
		Timeline: event.RoomEventFilter{
			Limit:           eventsToFetchPerRoom,
			LazyLoadMembers: true,
		},
	},
}

// Open starts the event loop of the client.
func (c *Client) Open() error {
	c.closeDone = make(chan struct{})
	filterID, err := c.FilterAdd(c.Filter)
	if err != nil {
		return err
	}
	go c.readLoop(filterID)

	return nil
}

// Close signals to the event loop to stop and wait for it to finish.
func (c *Client) Close() error {
	if c.shouldClose {
		return ErrAlreadyClosed
	}
	c.shouldClose = true
	<-c.closeDone

	return nil
}

func (c *Client) readLoop(filter string) {
	next := ""
	for !c.shouldClose {
		// Fetch next set of events.
		debug.Debug("Fetching new events. Next: " + next)
		resp, err := c.Sync(api.SyncArg{
			Filter:  filter,
			Since:   next,
			Timeout: syncTimeout,
		})
		if err != nil {
			// Exponentially backoff with a cap of 5 minutes.
			c.nextRetryTime *= 2
			if c.nextRetryTime < minBackoffTime {
				c.nextRetryTime = minBackoffTime
			}
			if c.nextRetryTime > maxBackoffTime {
				c.nextRetryTime = maxBackoffTime
			}

			debug.Error(fmt.Errorf("error in event loop (retrying in %d seconds): %w", c.nextRetryTime, err))
			time.Sleep(time.Duration(c.nextRetryTime) * time.Second)
			continue
		}

		handleWithRoomID := func(e []event.Event, roomID matrix.RoomID) {
			// Handle all events in the list.
			for _, v := range e {
				v.RoomID = roomID
				concrete, err := v.Parse()
				switch {
				case next == "":
					// Don't handle historical events.
					continue
				case errors.Is(err, event.ErrUnknownEventType):
					debug.Warn(fmt.Sprintf("unknown event type: %s", v.Type))
					continue
				case err != nil:
					debug.Warn(fmt.Errorf("error unmarshalling content: %w", err))
					continue
				}
				c.Handler.Handle(c, concrete)
			}
		}
		handle := func(e []event.Event) {
			handleWithRoomID(e, "")
		}

		handle(resp.Presence.Events)
		handle(resp.AccountData.Events)
		handle(resp.ToDevice.Events)
		for k, v := range resp.Rooms.Joined {
			handleWithRoomID(v.State.Events, k)
			handleWithRoomID(v.Timeline.Events, k)
			handleWithRoomID(v.Ephemeral.Events, k)
			handleWithRoomID(v.AccountData.Events, k)
		}
		// TODO resp.Rooms.Invited
		for k, v := range resp.Rooms.Left {
			handleWithRoomID(v.State.Events, k)
			handleWithRoomID(v.Timeline.Events, k)
			handleWithRoomID(v.AccountData.Events, k)
		}

		next = resp.NextBatch
	}

	c.closeDone <- struct{}{}
}
