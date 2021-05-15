package event

// Type is the type of the event that is contained in the contents field.
type Type string

// List of all known room events.
// NOTE: Update event/parse.go as well.
const (
	TypeRoomCanonicalAlias Type = "m.room.canonical_alias"
	TypeRoomCreate         Type = "m.room.create"
	TypeRoomJoinRules      Type = "m.room.join_rules"
	TypeRoomMember         Type = "m.room.member"
	TypeRoomPowerLevels    Type = "m.room.power_levels"
	TypeRoomRedaction      Type = "m.room.redaction"

	// Events from Instant Messaging module.
	TypeRoomMessage Type = "m.room.message"
	TypeRoomName    Type = "m.room.name"
	TypeRoomTopic   Type = "m.room.topic"
	TypeRoomAvatar  Type = "m.room.avatar"
	TypeRoomPinned  Type = "m.room.pinned_events"

	// Events from Voice over IP module.
	TypeCallInvite     Type = "m.call.invite"
	TypeCallCandidates Type = "m.call.candidates"
	TypeCallAnswer     Type = "m.call.answer"
	TypeCallHangup     Type = "m.call.hangup"

	// Events from Typing Notifications module.
	TypeTyping Type = "m.typing"

	// Events from Receipts module.
	TypeReceipt Type = "m.receipt"
)
