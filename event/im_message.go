package event

import (
	"encoding/json"
	"strings"

	"github.com/chanbakjsd/gotrix/encrypt"
	"github.com/chanbakjsd/gotrix/matrix"
)

var _ RoomEvent = RoomMessageEvent{}

// RoomMessageEvent represents a room event where a message has been sent.
//
// It has the type ID of `m.room.message`.
type RoomMessageEvent struct {
	RoomEventInfo `json:"-"`

	Body    string      `json:"body"`
	MsgType MessageType `json:"msgtype"`

	// This message is a reply to RelatesTo if present.
	RelatesTo json.RawMessage `json:"m.relates_to,omitempty"`

	// Optionally present in Text, Emote and Notice.
	Format        MessageFormat `json:"format,omitempty"`
	FormattedBody string        `json:"formatted_body,omitempty"`

	// This field is present in Location.
	GeoURI matrix.GeoURI `json:"geo_uri,omitempty"`

	// These fields are present in Image, File, Audio, Video.
	URL  matrix.URL    `json:"url,omitempty"`  // Present if content is not encrypted.
	File *encrypt.File `json:"file,omitempty"` // Present if content is encrypted.

	// This field is present in Image, File, Audio, Video, Location.
	// The relevant parsing functions should be used.
	Info json.RawMessage `json:"info,omitempty"` // Also present in Location.
}

// MessageType is the type of message sent.
type MessageType string

// All possible RoomMessageEvent types.
// List available at https://matrix.org/docs/spec/client_server/r0.6.1#m-room-message-msgtypes.
const (
	// Text, Emote and Notice are all messages.
	// Text is a regular message, Emote is similar to /me in IRC and Notice is a message sent by a bot.
	RoomMessageText   MessageType = "m.text"
	RoomMessageEmote  MessageType = "m.emote"
	RoomMessageNotice MessageType = "m.notice"

	RoomMessageImage    MessageType = "m.image"
	RoomMessageFile     MessageType = "m.file"
	RoomMessageAudio    MessageType = "m.audio"
	RoomMessageLocation MessageType = "m.location"
	RoomMessageVideo    MessageType = "m.video"
)

// MessageFormat is the type of the custom formatted body.
type MessageFormat string

// Currently, HTML is the only known RoomMessageFormat.
const (
	FormatHTML MessageFormat = "org.matrix.custom.html"
)

// ThumbnailInfo stores the info of a thumbnail.
type ThumbnailInfo struct {
	Height   int    `json:"h,omitempty"`        // Intended height of thumbnail.
	Width    int    `json:"w,omitempty"`        // Intended width of thumbnail.
	MimeType string `json:"mimetype,omitempty"` // MIME type of thumbnail.
	Size     int    `json:"size,omitempty"`     // Size in bytes.
}

// FileInfo stores the info of a file.
type FileInfo struct {
	MimeType      string        `json:"mimetype,omitempty"`       // MIME type of image.
	Size          int           `json:"size,omitempty"`           // Size in bytes.
	ThumbnailURL  matrix.URL    `json:"thumbnail_url,omitempty"`  // Present if thumbnail is unencrypted.
	ThumbnailFile encrypt.File  `json:"thumbnail_file,omitempty"` // Present if thumbnail is encrypted.
	ThumbnailInfo ThumbnailInfo `json:"thumbnail_info,omitempty"`
}

// ImageInfo stores the info of an image.
type ImageInfo struct {
	FileInfo

	// Intended display size of image. Present if RoomMessageFileInfo is part of RoomMessageImage.
	Height int `json:"h,omitempty"`
	Width  int `json:"w,omitempty"`
}

// AudioInfo stores the info of an audio.
type AudioInfo struct {
	Duration int    // Duration of audio in millisecond.
	MimeType string // MIME type of audio.
	Size     int    // Size in bytes.
}

// LocationInfo stores the info of a location.
type LocationInfo struct {
	ThumbnailURL  matrix.URL    `json:"thumbnail_url,omitempty"`  // Present if thumbnail is unencrypted.
	ThumbnailFile encrypt.File  `json:"thumbnail_file,omitempty"` // Present if thumbnail is encrypted.
	ThumbnailInfo ThumbnailInfo `json:"thumbnail_info,omitempty"`
}

// VideoInfo stores the info of a single video clip.
type VideoInfo struct {
	ImageInfo
	Duration int `json:"duration,omitempty"` // Duration of video in milliseconds.
}

// Type satisfies RoomEvent.
func (RoomMessageEvent) Type() Type {
	return TypeRoomMessage
}

// ImageInfo parses info as an ImageInfo.
func (e RoomMessageEvent) ImageInfo() (ImageInfo, error) {
	var a ImageInfo
	err := json.Unmarshal(e.Info, &a)
	return a, err
}

// FileInfo parses info as a FileInfo.
func (e RoomMessageEvent) FileInfo() (FileInfo, error) {
	var a FileInfo
	err := json.Unmarshal(e.Info, &a)
	return a, err
}

// AudioInfo parses info as an AudioInfo.
func (e RoomMessageEvent) AudioInfo() (AudioInfo, error) {
	var a AudioInfo
	err := json.Unmarshal(e.Info, &a)
	return a, err
}

// VideoInfo parses info as a VideoInfo.
func (e RoomMessageEvent) VideoInfo() (VideoInfo, error) {
	var a VideoInfo
	err := json.Unmarshal(e.Info, &a)
	return a, err
}

// LocationInfo parses info as a LocationInfo.
func (e RoomMessageEvent) LocationInfo() (LocationInfo, error) {
	var a LocationInfo
	err := json.Unmarshal(e.Info, &a)
	return a, err
}

// StrippedBody should be used if the client is rich reply aware (uses the RelatesTo field) in place
// of Body.
func (e RoomMessageEvent) StrippedBody() string {
	split := strings.Split(e.Body, "\n")
	var line int // Amount of line to strip.
	for line < len(split) && strings.HasPrefix(split[line], "> ") {
		line++
	}
	if line == len(split) {
		return ""
	}
	return strings.Join(split[line:], "\n")
}

// SetRoomEventInfo sets the room event info.
func (e *RoomMessageEvent) SetRoomEventInfo(i RoomEventInfo) {
	e.RoomEventInfo = i
}

// InReplyTo parses the message's RelatesTo object and returns the event ID that the message replies
// to, if any. An empty event string is returned if the message does not reply to another event in a
// valid way.
func (e *RoomMessageEvent) InReplyTo() matrix.EventID {
	var relatesTo struct {
		InReplyTo struct {
			EventID matrix.EventID `json:"event_id"`
		} `json:"m.in_reply_to"`
	}

	if err := json.Unmarshal(e.RelatesTo, &relatesTo); err != nil {
		return ""
	}

	return relatesTo.InReplyTo.EventID
}

// TODO Add helper method to parse RoomMessageHTML messages.
