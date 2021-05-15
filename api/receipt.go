package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/matrix"
)

// ReceiptType is the type of receipt. ReceiptRead is the only valid type.
type ReceiptType string

// ReceiptRead acknowledges that the event has been read.
const ReceiptRead ReceiptType = "m.read"

// ReceiptMarkerUpdate updates the location of receipt marker to the event ID specified.
func (c *Client) ReceiptMarkerUpdate(roomID matrix.RoomID, receiptType ReceiptType, eventID matrix.EventID) error {
	err := c.Request(
		"POST", EndpointRoomReceipt(roomID, receiptType, eventID), nil,
	)
	if err != nil {
		return fmt.Errorf("error updating receipt marker location: %w", err)
	}
	return nil
}
