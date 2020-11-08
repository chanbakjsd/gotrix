package event

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

// RoomNameEvent represents a state event where the room name is set.
// This is only used to be displayed. It's not unique and names can be duplicated.
//
// It has the type ID of `m.room.name` and a zero-length state key.
type RoomNameEvent struct {
	Name string `json:"name,omitempty"` // This must not exceed 255 bytes.
}

// RoomTopicEvent represents a state event where the room topic is set.
//
// It has the type ID of `m.room.topic` and a zero-length state key.
type RoomTopicEvent struct {
	Topic string `json:"topic,omitempty"`
}

// RoomAvatarEvent represents a state event where the room avatar is set.
//
// It has the type ID of `m.room.avatar` and a zero-length state key.
type RoomAvatarEvent struct {
	Info ImageInfo `json:"info,omitempty"`
	URL  string    `json:"url"`
}

// RoomPinnedEvent represents a state event where the list of events pinned are modified.
//
// It has the type ID of `m.room.pinned_events` and a zero-length state key.
type RoomPinnedEvent struct {
	Pinned []matrix.EventID `json:"pinned"`
}

// ContentOf implements EventContent.
func (e RoomNameEvent) ContentOf() Type {
	return TypeRoomName
}

// ContentOf implements EventContent.
func (e RoomTopicEvent) ContentOf() Type {
	return TypeRoomTopic
}

// ContentOf implements EventContent.
func (e RoomAvatarEvent) ContentOf() Type {
	return TypeRoomAvatar
}

// ContentOf implements EventContent.
func (e RoomPinnedEvent) ContentOf() Type {
	return TypeRoomPinned
}
