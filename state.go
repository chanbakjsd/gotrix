package gotrix

import (
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// State represents the required functions for a state.
type State interface {
	// RoomEvent returns the latest event in a room with the specified type.
	// If it is found in the cache, error will be nil.
	// Note that (nil, nil) should be returned if the cache can be certain the event type never occurred.
	RoomEvent(roomID matrix.RoomID, eventType event.Type, stateKey string) (e *event.Event, err error)
	// RoomEvents returns all the events with the given event type.
	// If there is duplicate events with the same state key, the newer one should be returned.
	RoomEvents(roomID matrix.RoomID, eventType event.Type) (map[string]*event.Event, error)
	// RoomEventSet creates a room event in the specified room.
	RoomEventSet(roomID matrix.RoomID, e *event.Event) error
}

// RoomState queries the internal State for the given RoomEvent.
// If the State does not have that event, it queries the homeserver directly.
func (c *Client) RoomState(roomID matrix.RoomID, eventType event.Type, key string) (*event.Event, error) {
	e, err := c.State.RoomEvent(roomID, eventType, key)
	if err == nil {
		return e, nil
	}

	return c.Client.RoomState(roomID, eventType, key)
}

// RoomStates queries the internal State for all state events that match RoomEvents.
func (c *Client) RoomStates(roomID matrix.RoomID, eventType event.Type) (map[string]*event.Event, error) {
	return c.State.RoomEvents(roomID, eventType)
}
