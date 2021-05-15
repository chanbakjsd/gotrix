package event

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

var _ StateEvent = RoomTombstoneEvent{}

// RoomTombstoneEvent is an event where the current room has been upgraded and a new room should be used instead.
type RoomTombstoneEvent struct {
	RoomEventInfo
	Message         string        `json:"body,omitempty"`
	ReplacementRoom matrix.RoomID `json:"replacement_room,omitempty"`
}

// Type implements StateEvent.
func (RoomTombstoneEvent) Type() Type {
	return TypeRoomTombstone
}

// StateKey implements StateEvent.
func (RoomTombstoneEvent) StateKey() string {
	return ""
}
