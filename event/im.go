package event

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

var (
	_ StateEvent = &RoomNameEvent{}
	_ StateEvent = &RoomTopicEvent{}
	_ StateEvent = &RoomAvatarEvent{}
	_ StateEvent = &RoomPinnedEvent{}
)

// RoomNameEvent represents a state event where the room name is set.
// This is only used to be displayed. It's not unique and names can be duplicated.
type RoomNameEvent struct {
	StateEventInfo `json:"-"`
	Name           string `json:"name,omitempty"` // This must not exceed 255 bytes.
}

// RoomTopicEvent represents a state event where the room topic is set.
type RoomTopicEvent struct {
	StateEventInfo `json:"-"`
	Topic          string `json:"topic,omitempty"`
}

// RoomAvatarEvent represents a state event where the room avatar is set.
type RoomAvatarEvent struct {
	StateEventInfo `json:"-"`
	Image          ImageInfo  `json:"info,omitempty"`
	URL            matrix.URL `json:"url"`
}

// RoomPinnedEvent represents a state event where the list of events pinned are modified.
type RoomPinnedEvent struct {
	StateEventInfo `json:"-"`
	Pinned         []matrix.EventID `json:"pinned"`
}
