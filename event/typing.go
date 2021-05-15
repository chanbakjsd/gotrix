package event

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

var _ Event = TypingEvent{}

// TypingEvent is an event that updates the list of users that are typing.
type TypingEvent struct {
	UserID []matrix.UserID `json:"user_ids"`
	RoomID matrix.RoomID   `json:"-"`
}

// Type implements Event.
func (TypingEvent) Type() Type {
	return TypeTyping
}
