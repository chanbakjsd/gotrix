package state

import (
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// RoomState is the state kept by a DefaultState for each room.
type RoomState map[event.Type]map[string]event.StateEvent

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

// RoomState returns the last event set by RoomEventSet.
// It never returns an error as it does not forget state.
func (d DefaultState) RoomState(roomID matrix.RoomID, eventType event.Type, key string) (event.StateEvent, error) {
	return d.TrackedState[roomID][eventType][key], nil
}

// RoomStates returns the last set of events set by RoomEventSet.
func (d DefaultState) RoomStates(roomID matrix.RoomID, eventType event.Type) (map[string]event.StateEvent, error) {
	return d.TrackedState[roomID][eventType], nil
}

// RoomStateSet sets the state inside a DefaultState to be returned by DefaultState later.
func (d *DefaultState) RoomStateSet(roomID matrix.RoomID, e event.StateEvent) error {
	eventType := e.Type()
	stateKey := e.StateKey()

	if _, ok := d.TrackedState[roomID]; !ok {
		d.TrackedState[roomID] = make(RoomState)
	}

	if _, ok := d.TrackedState[roomID][eventType]; !ok {
		d.TrackedState[roomID][eventType] = make(map[string]event.StateEvent)
	}

	d.TrackedState[roomID][eventType][stateKey] = e
	return nil
}
