package state

import (
	"sync"

	"github.com/chanbakjsd/gotrix/api"
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// ErrStopIter is a copy of gotrix.ErrStopIter.
var ErrStopIter error

// RoomState is the state kept by a DefaultState for each room.
type RoomState map[event.Type]map[string]event.StateEvent

// DefaultState is the default used implementation of state by the gotrix package.
type DefaultState struct {
	mu             sync.RWMutex
	roomStateMap   map[matrix.RoomID]RoomState
	roomSummaryMap map[matrix.RoomID]api.SyncRoomSummary
}

// NewDefault returns a DefaultState that has been initialized empty.
func NewDefault() *DefaultState {
	return &DefaultState{
		roomStateMap: make(map[matrix.RoomID]RoomState),
	}
}

// RoomState returns the last event added in AddEvents.
// It never returns an error as it does not forget state.
func (d *DefaultState) RoomState(roomID matrix.RoomID, eventType event.Type, key string) (event.StateEvent, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.roomStateMap[roomID][eventType][key], nil
}

// EachRoomState calls f for every event of the specified type.
// It terminates iteration when an error is returned. Use ErrStopIter to denote a non-failure condition.
// It makes a copy internally and calls f on it.
func (d *DefaultState) EachRoomState(id matrix.RoomID, typ event.Type, f func(string, event.StateEvent) error) error {
	states, err := d.RoomStates(id, typ)
	if err != nil {
		return err
	}

	for k, v := range states {
		err = f(k, v)
		switch {
		case err == ErrStopIter:
			return nil
		case err != nil:
			return err
		}
	}

	return nil
}

// RoomStates returns the last set of events added in AddEvents.
func (d *DefaultState) RoomStates(roomID matrix.RoomID, eventType event.Type) (map[string]event.StateEvent, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	events, ok := d.roomStateMap[roomID][eventType]
	if !ok {
		return nil, nil
	}

	eventsCopy := make(map[string]event.StateEvent, len(events))
	for k, v := range events {
		eventsCopy[k] = v
	}

	return eventsCopy, nil
}

// RoomSummary returns the summaries set in AddEvents.
func (d *DefaultState) RoomSummary(roomID matrix.RoomID) (api.SyncRoomSummary, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.roomSummaryMap[roomID], nil
}

func accumulateRaw(dst []event.StateEvent, roomID matrix.RoomID, raws []event.RawEvent) []event.StateEvent {
	for _, raw := range raws {
		e, err := event.Parse(raw)
		if err != nil {
			continue
		}
		state, ok := e.(event.StateEvent)
		if ok {
			state.RoomInfo().RoomID = roomID
			dst = append(dst, state)
		}
	}

	return dst
}

func accumulateStripped(dst []event.StateEvent, roomID matrix.RoomID, evs []event.StrippedEvent) []event.StateEvent {
	for _, ev := range evs {
		e, err := event.Parse(event.RawEvent(ev))
		if err != nil {
			continue
		}
		state, ok := e.(event.StateEvent)
		if ok {
			state.RoomInfo().RoomID = roomID
			dst = append(dst, state)
		}
	}

	return dst
}

// AddEvents sets the room state events inside a DefaultState to be returned by DefaultState later.
func (d *DefaultState) AddEvents(sync *api.SyncResponse) error {
	var eventCount int
	for _, v := range sync.Rooms.Joined {
		eventCount += len(v.State.Events)
	}
	for _, v := range sync.Rooms.Invited {
		eventCount += len(v.State.Events)
	}
	for _, v := range sync.Rooms.Left {
		eventCount += len(v.State.Events)
	}

	stateEvents := make([]event.StateEvent, 0, eventCount)
	for k, v := range sync.Rooms.Joined {
		stateEvents = accumulateRaw(stateEvents, k, v.State.Events)
	}
	for k, v := range sync.Rooms.Invited {
		stateEvents = accumulateStripped(stateEvents, k, v.State.Events)
	}
	for k, v := range sync.Rooms.Left {
		stateEvents = accumulateRaw(stateEvents, k, v.State.Events)
	}

	if len(stateEvents) == 0 {
		return nil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	for _, state := range stateEvents {
		info := state.StateInfo()
		roomID := info.RoomID
		stateKey := info.StateKey
		eventType := info.Type

		if _, ok := d.roomStateMap[roomID]; !ok {
			d.roomStateMap[roomID] = make(RoomState, 1)
		}

		if _, ok := d.roomStateMap[roomID][eventType]; !ok {
			d.roomStateMap[roomID][eventType] = make(map[string]event.StateEvent, 1)
		}

		d.roomStateMap[roomID][eventType][stateKey] = state
	}

	for k, v := range sync.Rooms.Joined {
		// Should always be larger than 0 as the user is in the room. If it is 0, it's empty.
		if v.Summary.JoinedCount > 0 {
			d.roomSummaryMap[k] = v.Summary
		}
	}

	return nil
}
