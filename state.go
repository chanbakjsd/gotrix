package gotrix

import (
	"errors"

	"github.com/chanbakjsd/gotrix/api"
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// ErrInvalidStateEvent is returned by (*Client).RoomState and (*Client).RoomStates when homeserver returns an event
// that is not a known state event when a state event is expected.
var ErrInvalidStateEvent = errors.New("invalid state event has been returned by homeserver")

// State represents the required functions for a state.
type State interface {
	// RoomState returns the latest event in a room with the specified type.
	// If it is found in the cache, error will be nil.
	// Note that (nil, nil) should be returned if the cache can be certain the event type never occurred.
	RoomState(roomID matrix.RoomID, eventType event.Type, stateKey string) (event.StateEvent, error)
	// RoomStates returns all the events with the given event type.
	// If there is duplicate events with the same state key, the newer one should be returned.
	RoomStates(roomID matrix.RoomID, eventType event.Type) (map[string]event.StateEvent, error)
	// AddEvent adds the needed events from the given sync response.
	// It is up to the implementation to pick and add the needed events inside the response.
	AddEvents(*api.SyncResponse) error
}

// RoomState queries the internal State for the given RoomEvent.
// If the State does not have that event, it queries the homeserver directly.
func (c *Client) RoomState(roomID matrix.RoomID, eventType event.Type, key string) (event.StateEvent, error) {
	e, err := c.State.RoomState(roomID, eventType, key)
	if err == nil {
		return e, nil
	}

	raw, err := c.Client.RoomState(roomID, eventType, key)
	if err != nil {
		return nil, err
	}

	parsed, err := raw.Parse()
	if err != nil {
		return nil, err
	}

	stateEvent, ok := parsed.(event.StateEvent)
	if !ok {
		return nil, ErrInvalidStateEvent
	}

	return stateEvent, nil
}

// RoomStates queries the internal State for all state events that match RoomEvents.
func (c *Client) RoomStates(roomID matrix.RoomID, eventType event.Type) (map[string]event.StateEvent, error) {
	return c.State.RoomStates(roomID, eventType)
}
