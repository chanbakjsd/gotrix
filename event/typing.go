package event

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

var _ Event = &TypingEvent{}

// TypingEvent is an event that updates the list of users that are typing.
type TypingEvent struct {
	EventInfo `json:"-"`

	UserID []matrix.UserID `json:"user_ids"`
	RoomID matrix.RoomID   `json:"room_id"`
}
