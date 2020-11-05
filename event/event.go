package event

import (
	"encoding/json"

	"github.com/chanbakjsd/gomatrix/matrix"
)

// StrippedEvent represents an event that has been stripped.
// This allows the client to display a room state correctly without its full timeline.
type StrippedEvent struct {
	Type     Type            `json:"type"`
	Content  json.RawMessage `json:"content"`
	StateKey string          `json:"state_key"`
	Sender   string          `json:"sender"`
}

// Event represents events that can be sent from homeserver to the client.
type Event struct {
	// Common data for all events.
	Type    Type            `json:"type"`
	Content json.RawMessage `json:"content"`

	// Data that are common for rooms and state events.
	EventID          string           `json:"event_id,omitempty"`
	Sender           matrix.UserID    `json:"sender,omitempty"`
	OriginServerTime matrix.Timestamp `json:"origin_server_ts,omitempty"`
	RoomID           string           `json:"room_id,omitempty"` // NOT included on `/sync` events.
	Unsigned         struct {
		// Age is the time in milliseconds that has elapsed since the event was sent.
		// It is generated by local homeserver and may be incorrect if either server's
		// time is out of sync.
		Age           matrix.Duration `json:"age,omitempty"`
		RedactReason  *Event          `json:"redacted_because,omitempty"`
		TransactionID string          `json:"transaction_id,omitempty"`
	} `json:"unsigned,omitempty"`

	// Data for state events.
	StateKey    string          `json:"state_key,omitempty"`
	PrevContent json.RawMessage `json:"prev_content"` // Optional previous content, if available.

	// Data for `m.room.redaction`. The ID of the event that was actually redacted.
	Redacts string `json:"redacts,omitempty"`
}
