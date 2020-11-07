package gomatrix

import (
	"time"

	"github.com/chanbakjsd/gomatrix/api"
	"github.com/chanbakjsd/gomatrix/debug"
	"github.com/chanbakjsd/gomatrix/event"
)

const (
	eventsToFetchPerRoom = 50
	minBackoffTime       = 1
	maxBackoffTime       = 300
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
		// TODO Already closed
		return nil
	}
	c.shouldClose = true
	<-c.closeDone

	return nil
}

func (c *Client) readLoop(filter string) {
	next := ""
	for !c.shouldClose {
		debug.Fields(map[string]interface{}{
			"next_id": next,
		}).Debug("Fetching new events.")
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

			debug.Fields(map[string]interface{}{
				"err":        err,
				"retry_time": c.nextRetryTime,
			}).Error("Event loop error!")

			time.Sleep(time.Duration(c.nextRetryTime) * time.Second)
			continue
		}

		next = resp.NextBatch
	}

	c.closeDone <- struct{}{}
}
