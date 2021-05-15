package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// RoomStateSendArg represents all required arguments to (*Client).RoomStateSend.
type RoomStateSendArg struct {
	Type     event.Type
	StateKey string
	Content  interface{}
}

// RoomStateSend sends the provided state event to the provided room ID.
func (c *Client) RoomStateSend(roomID matrix.RoomID, event RoomStateSendArg) (matrix.EventID, error) {
	var resp struct {
		EventID matrix.EventID `json:"event_id"`
	}
	err := c.Request(
		"PUT", EndpointRoomStateExact(roomID, event.Type, event.StateKey), &resp,
		httputil.WithToken(), httputil.WithJSONBody(event.Content),
	)
	if err != nil {
		return "", fmt.Errorf("error sending state event: %w", err)
	}
	return resp.EventID, nil
}

// RoomEventSend sends the provided one-off event to the provided room ID.
func (c *Client) RoomEventSend(roomID matrix.RoomID, eventType event.Type, body interface{}) (matrix.EventID, error) {
	var resp struct {
		EventID matrix.EventID `json:"event_id"`
	}

	err := c.Request(
		"PUT", EndpointRoomSend(roomID, eventType, NextTransactionID()), &resp,
		httputil.WithToken(), httputil.WithJSONBody(body),
	)
	if err != nil {
		return "", fmt.Errorf("error sending room event: %w", err)
	}
	return resp.EventID, nil
}

// RoomEventRedact redacts a room event as specified by the room ID and event ID.
// A user can redact events they sent out or other people's event provided they have the power level to.
func (c *Client) RoomEventRedact(roomID matrix.RoomID, eventID matrix.EventID, reason string) (matrix.EventID, error) {
	var resp struct {
		EventID matrix.EventID `json:"event_id"`
	}
	var err error

	if reason == "" {
		err = c.Request(
			"PUT", EndpointRoomRedact(roomID, eventID, NextTransactionID()), &resp,
			httputil.WithToken(),
		)
	} else {
		err = c.Request(
			"PUT", EndpointRoomRedact(roomID, eventID, NextTransactionID()), &resp,
			httputil.WithToken(), httputil.WithJSONBody(struct {
				Reason string `json:"reason"`
			}{
				Reason: reason,
			}))
	}
	if err != nil {
		return "", fmt.Errorf("error redacting room event: %w", err)
	}
	return resp.EventID, nil
}
