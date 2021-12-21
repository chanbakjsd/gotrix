package event

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

var _ Event = ReceiptEvent{}

// ReceiptEvent is an event where the read marker is updated.
type ReceiptEvent struct {
	*EventInfo

	Events map[matrix.EventID]Receipt `json:"content"`
	RoomID matrix.RoomID              `json:"room_id"`
}

// Receipt is an aggregate of users that have acknowledged a certain event.
type Receipt struct {
	Read map[matrix.UserID]struct {
		Timestamp int `json:"ts"`
	} `json:"m.read"`
}
