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

// SyncOptions contains options for the /sync endpoint that is used once the
// Client is opened.
type SyncOptions struct {
	Filter         event.GlobalFilter
	Timeout        time.Duration
	MinBackoffTime time.Duration
	MaxBackoffTime time.Duration
}

// DefaultSyncOptions is the default sync options instance used on every Client
// creation.
var DefaultSyncOptions = SyncOptions{
	Filter: event.GlobalFilter{
		Room: event.RoomFilter{
			IncludeLeave: false,
			State: event.StateFilter{
				LazyLoadMembers: true,
			},
			Timeline: event.RoomEventFilter{
				Limit:           50,
				LazyLoadMembers: true,
			},
		},
	},
	Timeout:        5 * time.Second,
	MinBackoffTime: 1 * time.Second,
	MaxBackoffTime: 300 * time.Second,
}

// Next returns the current Next synchronization argument. Next can ONLY be
// called once the Client is closed, otherwise a panic will occur.
func (c *Client) Next() string {
	select {
	case <-c.closeDone:
		return c.next
	default:
		panic("Next called on unclosed Client")
	}
}

// Open starts the event loop of the client with a background context.
func (c *Client) Open() error {
	return c.OpenWithNext("")
}

// syncOpts is the internal copy of the sync states.
type syncOpts struct {
	SyncOptions
	next     string
	filterID string
}

// OpenWithNext starts the event loop with the given next string that resumes the sync loop.
// If next is empty, then an initial sync will be done.
func (c *Client) OpenWithNext(next string) error {
	ctx, cancel := context.WithCancel(context.Background())

	c.closeDone = make(chan struct{})
	c.cancelFunc = cancel

	filterID, err := c.FilterAdd(c.SyncOpts.Filter)
	if err != nil {
		return err
	}

	go c.readLoop(ctx, syncOpts{
		SyncOptions: c.SyncOpts,
		next:        next,
		filterID:    filterID,
	})

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
			continue
		}
		c.Handler.Handle(c, concrete)
	}
}

func (c *Client) readLoop(ctx context.Context, opts syncOpts) {
	client := c.WithContext(ctx)

	timeout := int(opts.Timeout / time.Millisecond)
	next := opts.next

	handle := func(e []event.RawEvent) {
		c.handleWithRoomID(e, "", next == "")
	}

	var nextRetryTime time.Duration

	timer := time.NewTimer(0)
	defer timer.Stop()

	defer close(c.closeDone)

	<-timer.C

	for {
		// Fetch next set of events.
		debug.Debug("Fetching new events. Next: " + next)
		resp, err := client.Sync(api.SyncArg{
			Filter:  opts.filterID,
			Since:   next,
			Timeout: timeout,
		})
		if err != nil {
			if ctx.Err() != nil {
				// The context has finished.
				return
			}
			// Exponentially backoff with a cap of 5 minutes.
			nextRetryTime *= 2
			if nextRetryTime < opts.MinBackoffTime {
				nextRetryTime = opts.MinBackoffTime
			}
			if nextRetryTime > opts.MaxBackoffTime {
				nextRetryTime = opts.MaxBackoffTime
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

		if err := c.State.AddEvents(resp); err != nil {
			debug.Debug(fmt.Errorf("error adding sync events to state: %w", err))
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
