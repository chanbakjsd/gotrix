package event

import (
	"encoding/json"

	"github.com/chanbakjsd/gotrix/encrypt"
	"github.com/chanbakjsd/gotrix/matrix"
)

// RoomMessage represents a room event where a message has been sent.
//
// It has the type ID of `m.room.message`.
type RoomMessage struct {
	Event
	Body    string          `json:"body"`
	MsgType RoomMessageType `json:"msgtype"`

	// Optionally present in Text, Emote and Notice.
	Format        RoomMessageFormat `json:"format,omitempty"`
	FormattedBody string            `json:"formatted_body,omitempty"`

	// This field is present in Location.
	GeoURI matrix.GeoURI `json:"geo_uri,omitempty"`

	// These fields are present in Image, File, Audio, Video.
	URL  matrix.URL   `json:"url,omitempty"`  // Present if content is not encrypted.
	File encrypt.File `json:"file,omitempty"` // Present if content is encrypted.

	// This field is present in Image, File, Audio, Video, Location.
	// The relevant parsing functions should be used.
	Info json.RawMessage `json:"info,omitempty"` // Also present in Location.
}

// RoomMessageType is the type of message sent.
type RoomMessageType string

// All possible RoomMessage types.
// List available at https://matrix.org/docs/spec/client_server/r0.6.1#m-room-message-msgtypes.
const (
	// Text, Emote and Notice are all messages.
	// Text is a regular message, Emote is similar to /me in IRC and Notice is a message sent by a bot.
	RoomMessageText     RoomMessageType = "m.text"
	RoomMessageEmote    RoomMessageType = "m.emote"
	RoomMessageNotice   RoomMessageType = "m.notice"
	RoomMessageImage    RoomMessageType = "m.image"
	RoomMessageFile     RoomMessageType = "m.file"
	RoomMessageAudio    RoomMessageType = "m.audio"
	RoomMessageLocation RoomMessageType = "m.location"
	RoomMessageVideo    RoomMessageType = "m.video"
)

// RoomMessageFormat is the type of the custom formatted body.
type RoomMessageFormat string

// Currently, RoomMessageHTML is the only known RoomMessageFormat.
const (
	RoomMessageHTML RoomMessageFormat = "org.matrix.custom.html"
)

// RoomMessageThumbnailInfo stores the info of a thumbnail.
type RoomMessageThumbnailInfo struct {
	Height   int    `json:"h,omitempty"`        // Intended height of thumbnail.
	Width    int    `json:"w,omitempty"`        // Intended width of thumbnail.
	MimeType string `json:"mimetype,omitempty"` // MIME type of thumbnail.
	Size     int    `json:"size,omitempty"`     // Size in bytes.
}

// RoomMessageFileInfo stores the info of a file.
type RoomMessageFileInfo struct {
	MimeType      string                   `json:"mimetype,omitempty"`       // MIME type of image.
	Size          int                      `json:"size,omitempty"`           // Size in bytes.
	ThumbnailURL  matrix.URL               `json:"thumbnail_url,omitempty"`  // Present if thumbnail is unencrypted.
	ThumbnailFile encrypt.File             `json:"thumbnail_file,omitempty"` // Present if thumbnail is encrypted.
	ThumbnailInfo RoomMessageThumbnailInfo `json:"thumbnail_info,omitempty"`
}

// RoomMessageImageInfo stores the info of an image.
type RoomMessageImageInfo struct {
	RoomMessageFileInfo

	// Intended display size of image. Present if RoomMessageFileInfo is part of RoomMessageImage.
	Height int `json:"h,omitempty"`
	Width  int `json:"w,omitempty"`
}

// RoomMessageAudioInfo stores the info of an audio.
type RoomMessageAudioInfo struct {
	Duration int    // Duration of audio in millisecond.
	MimeType string // MIME type of audio.
	Size     int    // Size in bytes.
}

// RoomMessageLocationInfo stores the info of a location.
type RoomMessageLocationInfo struct {
	ThumbnailURL  matrix.URL               `json:"thumbnail_url,omitempty"`  // Present if thumbnail is unencrypted.
	ThumbnailFile encrypt.File             `json:"thumbnail_file,omitempty"` // Present if thumbnail is encrypted.
	ThumbnailInfo RoomMessageThumbnailInfo `json:"thumbnail_info,omitempty"`
}

// RoomMessageVideoInfo stores the info of a single video clip.
type RoomMessageVideoInfo struct {
	RoomMessageImageInfo
	Duration int `json:"duration,omitempty"` // Duration of video in milliseconds.
}

// ContentOf implements EventContent.
func (e RoomMessage) ContentOf() Type {
	return TypeRoomMessage
}

// ImageInfo parses info as a RoomMessageImageInfo.
func (e RoomMessage) ImageInfo() (RoomMessageImageInfo, error) {
	var a RoomMessageImageInfo
	err := json.Unmarshal(e.Info, &a)
	return a, err
}

// FileInfo parses info as a RoomMessageFileInfo.
func (e RoomMessage) FileInfo() (RoomMessageFileInfo, error) {
	var a RoomMessageFileInfo
	err := json.Unmarshal(e.Info, &a)
	return a, err
}

// AudioInfo parses info as a RoomMessageAudioInfo.
func (e RoomMessage) AudioInfo() (RoomMessageAudioInfo, error) {
	var a RoomMessageAudioInfo
	err := json.Unmarshal(e.Info, &a)
	return a, err
}

// VideoInfo parses info as a RoomMessageImageInfo.
func (e RoomMessage) VideoInfo() (RoomMessageVideoInfo, error) {
	var a RoomMessageVideoInfo
	err := json.Unmarshal(e.Info, &a)
	return a, err
}

// LocationInfo parses info as a RoomMessageLocationInfo.
func (e RoomMessage) LocationInfo() (RoomMessageLocationInfo, error) {
	var a RoomMessageLocationInfo
	err := json.Unmarshal(e.Info, &a)
	return a, err
}

// TODO Add helper method to parse RoomMessageHTML messages.
