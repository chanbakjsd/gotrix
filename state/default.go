package state

import (
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// RoomState is the state kept by a DefaultState for each room.
type RoomState map[event.Type]map[string]*event.Event

// DefaultState is the default used implementation of state by the gotrix package.
type DefaultState struct {
	TrackedState map[matrix.RoomID]RoomState
}

// NewDefault returns a DefaultState that has been initialized empty.
func NewDefault() *DefaultState {
	return &DefaultState{
		TrackedState: make(map[matrix.RoomID]RoomState),
	}
}

// RoomEvent returns the last event set by RoomEventSet.
// It never returns an error as it does not forget state.
func (d DefaultState) RoomEvent(roomID matrix.RoomID, eventType event.Type, key string) (*event.Event, error) {
	return d.TrackedState[roomID][eventType][key], nil
}

// RoomEvents returns the last set of events set by RoomEventSet.
func (d DefaultState) RoomEvents(roomID matrix.RoomID, eventType event.Type) (map[string]*event.Event, error) {
	return d.TrackedState[roomID][eventType], nil
}

// RoomEventSet sets the state inside a DefaultState to be returned by DefaultState later.
func (d *DefaultState) RoomEventSet(roomID matrix.RoomID, e *event.Event) error {
	if _, ok := d.TrackedState[roomID]; !ok {
		d.TrackedState[roomID] = make(RoomState)
	}

	if _, ok := d.TrackedState[roomID][e.Type]; !ok {
		d.TrackedState[roomID][e.Type] = make(map[string]*event.Event)
	}

	d.TrackedState[roomID][e.Type][e.StateKey] = e
	return nil
}
