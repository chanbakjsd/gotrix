package event

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

var (
	_ StateEvent = RoomNameEvent{}
	_ StateEvent = RoomTopicEvent{}
	_ StateEvent = RoomAvatarEvent{}
	_ StateEvent = RoomPinnedEvent{}
)

// RoomNameEvent represents a state event where the room name is set.
// This is only used to be displayed. It's not unique and names can be duplicated.
//
// It has the type ID of `m.room.name` and a zero-length state key.
type RoomNameEvent struct {
	RoomEventInfo
	Name string `json:"name,omitempty"` // This must not exceed 255 bytes.
}

// RoomTopicEvent represents a state event where the room topic is set.
//
// It has the type ID of `m.room.topic` and a zero-length state key.
type RoomTopicEvent struct {
	RoomEventInfo
	Topic string `json:"topic,omitempty"`
}

// RoomAvatarEvent represents a state event where the room avatar is set.
//
// It has the type ID of `m.room.avatar` and a zero-length state key.
type RoomAvatarEvent struct {
	RoomEventInfo
	Info ImageInfo  `json:"info,omitempty"`
	URL  matrix.URL `json:"url"`
}

// RoomPinnedEvent represents a state event where the list of events pinned are modified.
//
// It has the type ID of `m.room.pinned_events` and a zero-length state key.
type RoomPinnedEvent struct {
	RoomEventInfo
	Pinned []matrix.EventID `json:"pinned"`
}

// Type satisfies StateEvent.
func (RoomNameEvent) Type() Type {
	return TypeRoomName
}

// Type satisfies StateEvent.
func (RoomTopicEvent) Type() Type {
	return TypeRoomTopic
}

// Type satisfies StateEvent.
func (RoomAvatarEvent) Type() Type {
	return TypeRoomAvatar
}

// Type satisfies StateEvent.
func (RoomPinnedEvent) Type() Type {
	return TypeRoomPinned
}

// StateKey satisfies StateEvent.
func (RoomNameEvent) StateKey() string {
	return ""
}

// StateKey satisfies StateEvent.
func (RoomTopicEvent) StateKey() string {
	return ""
}

// StateKey satisfies StateEvent.
func (RoomAvatarEvent) StateKey() string {
	return ""
}

// StateKey satisfies StateEvent.
func (RoomPinnedEvent) StateKey() string {
	return ""
}
