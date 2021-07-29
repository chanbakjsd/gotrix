package gotrix

import (
	"errors"

	"github.com/chanbakjsd/gotrix/api"
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
	"github.com/chanbakjsd/gotrix/state"
)

// ErrInvalidStateEvent is returned by (*Client).RoomState and (*Client).RoomStates when homeserver returns an event
// that is not a known state event when a state event is expected.
var ErrInvalidStateEvent = errors.New("invalid state event has been returned by homeserver")

// ErrStopIter is an error used to denote that the iteration on EachRoomState should be stopped.
var ErrStopIter = errors.New("stop iterating on EachRoomState")

func init() {
	state.ErrStopIter = ErrStopIter
}

// State represents the required functions for a state.
type State interface {
	// RoomState returns the latest event in a room with the specified type.
	// If it is found in the cache, error will be nil.
	// Note that (nil, nil) should be returned if the cache can be certain the event type never occurred.
	RoomState(roomID matrix.RoomID, eventType event.Type, stateKey string) (event.StateEvent, error)
	// EachRoomState calls f for every event stored in the state.
	// To abort iteration, f should return ErrStopIter.
	// This function can return the error returned by f or errors while getting data for iteration.
	EachRoomState(roomID matrix.RoomID, eventType event.Type, f func(key string, e event.StateEvent) error) error
	// RoomSummary returns the summary of a room as received in sync response.
	RoomSummary(roomID matrix.RoomID) (api.SyncRoomSummary, error)
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

// EachRoomState iterates through all events with the specified type, stopping if f returns ErrIterStop.
func (c *Client) EachRoomState(roomID matrix.RoomID, typ event.Type, f func(string, event.StateEvent) error) error {
	return c.State.EachRoomState(roomID, typ, f)
}

// RoomSummary queries the State for the summary of a room, commonly used for generating room name.
// To follow the room name generation strategy of the specification, use Client.RoomName instead.
func (c *Client) RoomSummary(roomID matrix.RoomID) (api.SyncRoomSummary, error) {
	return c.State.RoomSummary(roomID)
}
