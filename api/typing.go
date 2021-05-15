package api

import (
	"fmt"
	"time"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// TypingStart notifies the server that the user is typing in a specific room.
// It should be repeated while the user is typing with preferably a few seconds of safety margin from timeout.
func (c *Client) TypingStart(roomID matrix.RoomID, userID matrix.UserID, timeout time.Duration) error {
	req := map[string]interface{}{
		"typing":  true,
		"timeout": timeout / time.Millisecond,
	}

	err := c.Request(
		"PUT", EndpointRoomTyping(roomID, userID), nil,
		httputil.WithToken(), httputil.WithJSONBody(req),
	)
	if err != nil {
		return fmt.Errorf("error sending typing start notification: %w", err)
	}
	return nil
}

// TypingStop notifies the server that the user has stopped typing in a specific room.
func (c *Client) TypingStop(roomID matrix.RoomID, userID matrix.UserID) error {
	req := map[string]interface{}{
		"typing": false,
	}

	err := c.Request(
		"PUT", EndpointRoomTyping(roomID, userID), nil,
		httputil.WithToken(), httputil.WithJSONBody(req),
	)
	if err != nil {
		return fmt.Errorf("error sending typing stop notification: %w", err)
	}
	return nil
}
