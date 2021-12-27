package gotrix

import (
	"encoding/json"
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

// File is a file that can be uploaded to Matrix homeserver.
type File struct {
	Name     string
	MIMEType string
	Content  io.ReadCloser
	Caption  string // Set to filename if empty.

	FileInfo  *event.FileInfo
	ImageInfo *event.ImageInfo
	AudioInfo *event.AudioInfo
	VideoInfo *event.VideoInfo
}

// SendImage uploads the provided image to the server and sends a message containing it to the designated room.
func (c *Client) SendImage(roomID matrix.RoomID, file File) (matrix.EventID, error) {
	return c.sendFile(roomID, event.RoomMessageImage, file)
}

// SendFile uploads the provided file to the server and sends a message containing it to the designated room.
func (c *Client) SendFile(roomID matrix.RoomID, file File) (matrix.EventID, error) {
	return c.sendFile(roomID, event.RoomMessageFile, file)
}

// SendAudio uploads the provided audio file to the server and sends a message containing it to the designated room.
func (c *Client) SendAudio(roomID matrix.RoomID, file File) (matrix.EventID, error) {
	return c.sendFile(roomID, event.RoomMessageAudio, file)
}

// SendLocation sends the provided location to the provided room ID.
func (c *Client) SendLocation(roomID matrix.RoomID, geoURI matrix.GeoURI, caption string) (matrix.EventID, error) {
	return c.RoomEventSend(roomID, event.TypeRoomMessage, event.RoomMessageEvent{
		MessageType: event.RoomMessageLocation,
		Body:        caption,
		GeoURI:      geoURI,
	})
}

// SendVideo uploads the provided video file to the server and sends a message containing it to the designated room.
func (c *Client) SendVideo(roomID matrix.RoomID, file File) (matrix.EventID, error) {
	return c.sendFile(roomID, event.RoomMessageVideo, file)
}

func (c *Client) sendMessage(roomID matrix.RoomID, msgType event.MessageType, content string) (matrix.EventID, error) {
	return c.RoomEventSend(roomID, event.TypeRoomMessage, event.RoomMessageEvent{
		MessageType: msgType,
		Body:        content,
	})
}

// sendFile sends the image to the provided room ID with the provided content.
func (c *Client) sendFile(roomID matrix.RoomID, msgType event.MessageType, file File) (matrix.EventID, error) {
	url, err := c.MediaUpload(file.MIMEType, file.Name, file.Content)
	if err != nil {
		return "", err
	}

	if file.Caption == "" {
		file.Caption = file.Name
	}

	var additionalInfo json.RawMessage
	if file.AudioInfo != nil {
		additionalInfo, err = json.Marshal(file.AudioInfo)
		if err != nil {
			return "", err
		}
	}
	if file.FileInfo != nil {
		additionalInfo, err = json.Marshal(file.FileInfo)
		if err != nil {
			return "", err
		}
	}
	if file.ImageInfo != nil {
		additionalInfo, err = json.Marshal(file.ImageInfo)
		if err != nil {
			return "", err
		}
	}
	if file.VideoInfo != nil {
		additionalInfo, err = json.Marshal(file.VideoInfo)
		if err != nil {
			return "", err
		}
	}

	return c.RoomEventSend(roomID, event.TypeRoomMessage, event.RoomMessageEvent{
		MessageType:    msgType,
		Body:           file.Caption,
		URL:            url,
		AdditionalInfo: additionalInfo,
	})
}
