package event

// RoomMessage represents a room event where a message has been sent.
//
// It has the type ID of `m.room.message`.
type RoomMessage struct {
	Event
	Body    string          `json:"body"`
	MsgType RoomMessageType `json:"msgtype"`

	// TODO Add more fields
}

// RoomMessageType is the type of message sent.
type RoomMessageType string

// All possible RoomMessage types.
// List available at https://matrix.org/docs/spec/client_server/r0.6.1#m-room-message-msgtypes.
const (
	RoomMessageText     RoomMessageType = "m.text"
	RoomMessageNotice   RoomMessageType = "m.notice"
	RoomMessageImage    RoomMessageType = "m.image"
	RoomMessageFile     RoomMessageType = "m.file"
	RoomMessageAudio    RoomMessageType = "m.audio"
	RoomMessageLocation RoomMessageType = "m.location"
	RoomMessageVideo    RoomMessageType = "m.video"
)

// ContentOf implements EventContent.
func (e RoomMessage) ContentOf() Type {
	return TypeRoomMessage
}
