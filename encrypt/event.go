package encrypt

import (
	"encoding/json"

	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// Actual Matrix event

// List of actual Matrix events that is related to end-to-end encryption.
const (
	TypeRoomEncryption   event.Type = "m.room.encryption"
	TypeRoomEncrypted    event.Type = "m.room.encrypted"
	TypeRoomKey          event.Type = "m.room_key"
	TypeRoomKeyRequest   event.Type = "m.room_key_request"
	TypeForwardedRoomKey event.Type = "m.forwarded_room_key"
	TypeDummy            event.Type = "m.dummy"
	TypeRoomKeyWithheld  event.Type = "m.room_key.withheld"
)

var (
	_ event.StateEvent = &RoomEncryptionEvent{}
	_ event.RoomEvent  = &RoomEncryptedEvent{}
	_ event.Event      = &RoomKeyEvent{}
	_ event.Event      = &RoomKeyRequestEvent{}
	_ event.Event      = &ForwardedRoomKeyEvent{}
	_ event.Event      = &DummyEvent{}
	_ event.Event      = &RoomKeyWithheldEvent{}
)

// RoomEncryptionEvent is a state event that defines how messages in a room
// should be encrypted.
type RoomEncryptionEvent struct {
	event.StateEventInfo `json:"-"`

	Algorithm Algorithm `json:"algorithm"`

	// How often to rotate keys in milliseconds.
	RotationTime int `json:"rotation_period_ms"`
	// How often to rotate keys based on messages sent.
	RotationMessage int `json:"rotation_period_msgs"`
}

// Algorithm is the algorithm used for encryption.
type Algorithm string

// List of algorithms specified in the specification.
const (
	AlgorithmOlm    Algorithm = "m.olm.v1.curve25519-aes-sha2"
	AlgorithmMegOlm Algorithm = "m.megolm.v1.aes-sha2"
)

// RoomEncryptedEvent is an event that may be a room event or a send-to-device
// event. It contains encrypted data which should be decrypted and processed.
type RoomEncryptedEvent struct {
	event.RoomEventInfo `json:"-"`

	Algorithm  Algorithm       `json:"algorithm"`
	Ciphertext json.RawMessage `json:"ciphertext"`
	DeviceID   matrix.DeviceID `json:"device_id,omitempty"`
	SenderKey  Curve25519Key   `json:"sender_key"`
	SessionID  SessionID       `json:"session_id,omitempty"`
}

// RoomKeyEvent is an event used to send room keys, it is typically stored
// in a RoomEncryptedEvent.
type RoomKeyEvent struct {
	event.EventInfo `json:"-"`

	Algorithm  Algorithm     `json:"algorithm"`
	RoomID     matrix.RoomID `json:"room_id"`
	SessionID  SessionID     `json:"session_id"`
	SessionKey Key           `json:"session_key"`
}

// RoomKeyRequestEvent is an event used to request room keys, it is typically
// sent as an un-encrypted send-to-device event.
type RoomKeyRequestEvent struct {
	event.EventInfo `json:"-"`

	Action           RoomKeyRequestAction `json:"action"`
	Body             *RoomKeyRequestInfo  `json:"body,omitempty"`
	RequestID        string               `json:"request_id"`
	RequestingDevice matrix.DeviceID      `json:"requesting_device_id"`
}

// RoomKeyRequestAction is the action the RoomKeyRequestEvent is sent for.
type RoomKeyRequestAction string

// Different types of RoomKeyRequestAction.
const (
	RoomKeyRequestRequest RoomKeyRequestAction = "request"
	RoomKeyRequestCancel  RoomKeyRequestAction = "request_cancellation"
)

// RoomKeyRequestInfo is the info of the key being requested.
type RoomKeyRequestInfo struct {
	Algorithm Algorithm     `json:"algorithm"`
	RoomID    matrix.RoomID `json:"room_id"`
	SenderKey Curve25519Key `json:"session_key"`
	SessionID SessionID     `json:"session_id"`
}

// ForwardedRoomKeyEvent is used to forward keys for end-to-end encryption.
// Typically it is encrypted as a RoomEncryptedEvent, then sent as a to-device
// event.
type ForwardedRoomKeyEvent struct {
	event.EventInfo `json:"-"`

	Algorithm Algorithm `json:"algorithm"`
	// Chain of keys this room key was forwarded from.
	ForwardChain     []Curve25519Key `json:"forwarding_curve25519_key_chain"`
	RoomID           matrix.RoomID   `json:"room_id"`
	SenderClaimedKey Ed25519Key      `json:"sender_claimed_ed25519_key"`
	SenderKey        Curve25519Key   `json:"sender_key"`
	SessionID        SessionID       `json:"session_id"`
	SessionKey       Curve25519Key   `json:"session_key"`
	// Fields of the withheld event, but not the event itself.
	Withheld *RoomKeyWithheldEvent `json:"withheld,omitempty"`
}

// RoomKeyWithheldEvent is an event to notify that room keys have been
// withheld. Algorithm, code and sender key are always present in the event and
// one of room ID or session ID is present in the event depending on the
// RoomKeyWithheldReason.
type RoomKeyWithheldEvent struct {
	event.EventInfo `json:"-"`

	Algorithm Algorithm             `json:"algorithm,omitempty"`
	Code      RoomKeyWithheldReason `json:"code"`
	// Human-readable version of Code.
	Reason    string        `json:"reason,omitempty"`
	RoomID    matrix.RoomID `json:"room_id,omitempty"`
	SessionID SessionID     `json:"session_id,omitempty"`
	SenderKey Curve25519Key `json:"sender_key,omitempty"`
}

// RoomKeyWithheldReason is the reason the room key is withheld instead of sent.
type RoomKeyWithheldReason string

const (
	// WithheldReasonBlacklist withholds the room keys because the device
	// requesting it was blacklisted.
	WithheldReasonBlacklist RoomKeyWithheldReason = "m.blacklisted"
	// WithheldReasonUnverified withholds the room keys because the device
	// requesting it has not been verified yet.
	WithheldReasonUnverified RoomKeyWithheldReason = "m.unverified"
	// WithheldReasonUnauthorized withholds the room keys because the device
	// requesting it does not have permission to do so such as the user not
	// being in the room.
	WithheldReasonUnauthorized RoomKeyWithheldReason = "m.unauthorised"
	// WithheldReasonUnavailable withholds the room keys because the device
	// does not have the requested room keys.
	WithheldReasonUnavailable RoomKeyWithheldReason = "m.unavailable"
	// WithheldReasonOlmFailure withholds the room keys because an Olm session
	// cannot be established.
	WithheldReasonOlmFailure RoomKeyWithheldReason = "m.no_olm"
)

// DummyEvent is a dummy event typically used to initiate renegotiation of keys.
type DummyEvent struct {
	event.EventInfo `json:"-"`
}
