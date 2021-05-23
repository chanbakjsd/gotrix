package gotrix

import (
	"context"
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
	minBackoffTime       = 1 * time.Second
	maxBackoffTime       = 300 * time.Second
	syncTimeout          = 5000
)

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

// Open starts the event loop of the client with a background context.
func (c *Client) Open() error {
	return c.OpenCtx(context.Background())
}

// OpenCtx starts the event loop of the client with the provided context.
func (c *Client) OpenCtx(ctx context.Context) error {
	c.closeDone = make(chan struct{})
	ctx, c.cancelFunc = context.WithCancel(ctx)

	filterID, err := c.FilterAdd(c.Filter)
	if err != nil {
		return err
	}
	go c.readLoop(ctx, filterID)

	return nil
}

// Close signals to the event loop to stop and wait for it to finish.
func (c *Client) Close() error {
	c.cancelFunc()
	<-c.closeDone

	return nil
}

func (c *Client) handleWithRoomID(e []event.RawEvent, roomID matrix.RoomID, isHistorical bool) {
	for _, v := range e {
		v := v
		v.RoomID = roomID
		concrete, err := v.Parse()

		// Update state if it's a state event.
		if stateEvent, ok := concrete.(event.StateEvent); ok {
			_ = c.State.RoomStateSet(roomID, stateEvent)
		}

		// Print out warnings.
		switch {
		case errors.Is(err, event.ErrUnknownEventType):
			debug.Warn(fmt.Sprintf("unknown event type: %s", v.Type))
		case err != nil:
			debug.Warn(fmt.Errorf("error unmarshalling content: %w", err))
		}

		// Don't call handlers on historical events.
		if isHistorical {
			continue
		}

		c.Handler.HandleRaw(c, v)
		if err != nil {
			return
		}
		c.Handler.Handle(c, concrete)
	}
}

func (c *Client) readLoop(ctx context.Context, filter string) {
	client := c.WithContext(ctx)
	next := ""

	handle := func(e []event.RawEvent) {
		c.handleWithRoomID(e, "", next == "")
	}

	var nextRetryTime time.Duration

	timer := time.NewTimer(0)
	defer timer.Stop()

	<-timer.C

	for {
		// Fetch next set of events.
		debug.Debug("Fetching new events. Next: " + next)
		resp, err := client.Sync(api.SyncArg{
			Filter:  filter,
			Since:   next,
			Timeout: syncTimeout,
		})
		if err != nil {
			if ctx.Err() != nil {
				// The context has finished.
				close(c.closeDone)
				return
			}
			// Exponentially backoff with a cap of 5 minutes.
			nextRetryTime *= 2
			if nextRetryTime < minBackoffTime {
				nextRetryTime = minBackoffTime
			}
			if nextRetryTime > maxBackoffTime {
				nextRetryTime = maxBackoffTime
			}

			debug.Error(fmt.Errorf("error in event loop (retrying in %s): %w", nextRetryTime, err))
			timer.Reset(nextRetryTime)
			select {
			case <-timer.C:
				continue
			case <-ctx.Done():
				return
			}
		}

		handle(resp.Presence.Events)
		handle(resp.AccountData.Events)
		handle(resp.ToDevice.Events)
		for k, v := range resp.Rooms.Joined {
			c.handleWithRoomID(v.State.Events, k, next == "")
			c.handleWithRoomID(v.Timeline.Events, k, next == "")
			c.handleWithRoomID(v.Ephemeral.Events, k, next == "")
			c.handleWithRoomID(v.AccountData.Events, k, next == "")
		}
		for k, v := range resp.Rooms.Invited {
			events := make([]event.RawEvent, len(v.State.Events))
			for k, v := range v.State.Events {
				events[k] = v.RawEvent
			}
			c.handleWithRoomID(events, k, next == "")
		}
		for k, v := range resp.Rooms.Left {
			c.handleWithRoomID(v.State.Events, k, next == "")
			c.handleWithRoomID(v.Timeline.Events, k, next == "")
			c.handleWithRoomID(v.AccountData.Events, k, next == "")
		}

		next = resp.NextBatch
	}
}
