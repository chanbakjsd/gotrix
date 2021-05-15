package gotrix

import (
	"io"

	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// SendMessage sends a message to the provided room ID with the provided content.
func (c *Client) SendMessage(roomID matrix.RoomID, content string) (matrix.EventID, error) {
	return c.sendMessage(roomID, event.RoomMessageText, content)
}

// SendEmote sends a emote to the provided room ID with the provided content.
//
// Emote is a regular message but is sent as someone performing it (/me in IRC).
func (c *Client) SendEmote(roomID matrix.RoomID, content string) (matrix.EventID, error) {
	return c.sendMessage(roomID, event.RoomMessageEmote, content)
}

// SendNotice sends a notice to the provided room ID with the provided content.
//
// Notice are the same as messages except they're not intended to be parsed by bots (ie. other bots'
// messages).
func (c *Client) SendNotice(roomID matrix.RoomID, content string) (matrix.EventID, error) {
	return c.sendMessage(roomID, event.RoomMessageNotice, content)
}

// SendImage uploads the provided image to the server and sends a message containing it to the designated room.
func (c *Client) SendImage(roomID matrix.RoomID, mime string, name string, file io.ReadCloser,
	caption string) (matrix.EventID, error) {
	return c.sendFile(roomID, event.RoomMessageImage, mime, name, file, caption)
}

// SendFile uploads the provided file to the server and sends a message containing it to the designated room.
func (c *Client) SendFile(roomID matrix.RoomID, mime string, name string, file io.ReadCloser) (matrix.EventID, error) {
	return c.sendFile(roomID, event.RoomMessageFile, mime, name, file, name)
}

// SendAudio uploads the provided audio file to the server and sends a message containing it to the designated room.
func (c *Client) SendAudio(roomID matrix.RoomID, mime string, name string, file io.ReadCloser,
	caption string) (matrix.EventID, error) {
	return c.sendFile(roomID, event.RoomMessageAudio, mime, name, file, caption)
}

// SendLocation sends the provided location to the provided room ID.
func (c *Client) SendLocation(roomID matrix.RoomID, geoURI matrix.GeoURI, caption string) (matrix.EventID, error) {
	return c.RoomEventSend(roomID, event.TypeRoomMessage, event.RoomMessageEvent{
		MsgType: event.RoomMessageLocation,
		Body:    caption,
		GeoURI:  geoURI,
	})
}

// SendVideo uploads the provided video file to the server and sends a message containing it to the designated room.
func (c *Client) SendVideo(roomID matrix.RoomID, mime string, name string, file io.ReadCloser,
	caption string) (matrix.EventID, error) {
	return c.sendFile(roomID, event.RoomMessageVideo, mime, name, file, caption)
}

func (c *Client) sendMessage(roomID matrix.RoomID, msgType event.MessageType, content string) (matrix.EventID, error) {
	return c.RoomEventSend(roomID, event.TypeRoomMessage, event.RoomMessageEvent{
		MsgType: msgType,
		Body:    content,
	})
}

// sendFile sends the image to the provided room ID with the provided content.
func (c *Client) sendFile(roomID matrix.RoomID, msgType event.MessageType, fileMime string, fileName string,
	file io.ReadCloser, caption string) (matrix.EventID, error) {
	url, err := c.MediaUpload(fileMime, fileName, file)
	if err != nil {
		return "", err
	}

	return c.RoomEventSend(roomID, event.TypeRoomMessage, event.RoomMessageEvent{
		MsgType: msgType,
		Body:    caption,
		URL:     url,
	})
}
