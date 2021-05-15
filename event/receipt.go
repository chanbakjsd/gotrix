package event

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

var _ Event = ReceiptEvent{}

// ReceiptEvent is an event where the read marker is updated.
type ReceiptEvent struct {
	Events map[matrix.EventID]Receipt
	RoomID matrix.RoomID
}

// Receipt is an aggregate of users that have acknowledged a certain event.
type Receipt struct {
	Read map[matrix.UserID]struct {
		Timestamp int `json:"ts"`
	} `json:"m.read"`
}

// Type implements Event.
func (ReceiptEvent) Type() Type {
	return TypeReceipt
}
