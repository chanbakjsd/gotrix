package gotrix

import (
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// SendMessage sends a message to the provided room ID with the provided content.
func (c *Client) SendMessage(roomID matrix.RoomID, content string) (matrix.EventID, error) {
	return c.RoomEventSend(roomID, event.TypeRoomMessage, event.RoomMessage{
		MsgType: event.RoomMessageText,
		Body:    content,
	})
}

// SendNotice sends a notice to the provided room ID with the provided content.
//
// Notice are the same as messages except they're not intended to be parsed by bots (ie. other bots'
// messages).
func (c *Client) SendNotice(roomID matrix.RoomID, content string) (matrix.EventID, error) {
	return c.RoomEventSend(roomID, event.TypeRoomMessage, event.RoomMessage{
		MsgType: event.RoomMessageText,
		Body:    content,
	})
}
