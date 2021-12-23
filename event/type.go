package event

import "encoding/json"

// Type is the type of the event that is contained in the contents field.
type Type string

// List of all known room events.
// NOTE: Update the 'parser' variable below as well.
const (
	TypeRoomCanonicalAlias Type = "m.room.canonical_alias"
	TypeRoomCreate         Type = "m.room.create"
	TypeRoomJoinRules      Type = "m.room.join_rules"
	TypeRoomMember         Type = "m.room.member"
	TypeRoomPowerLevels    Type = "m.room.power_levels"
	TypeRoomRedaction      Type = "m.room.redaction"

	// Events from the Instant Messaging module.
	TypeRoomMessage Type = "m.room.message"
	TypeRoomName    Type = "m.room.name"
	TypeRoomTopic   Type = "m.room.topic"
	TypeRoomAvatar  Type = "m.room.avatar"
	TypeRoomPinned  Type = "m.room.pinned_events"

	// Events from the Direct Messaging module.
	TypeDirect Type = "m.direct"

	// Events from the Voice over IP module.
	TypeCallInvite     Type = "m.call.invite"
	TypeCallCandidates Type = "m.call.candidates"
	TypeCallAnswer     Type = "m.call.answer"
	TypeCallHangup     Type = "m.call.hangup"

	// Events from the Typing Notifications module.
	TypeTyping Type = "m.typing"

	// Events from the Receipts module.
	TypeReceipt Type = "m.receipt"

	// Events from the Presence module.
	TypePresence Type = "m.presence"

	// Events from the History Visibility module.
	TypeRoomHistoryVisibility Type = "m.room.history_visibility"

	// Events from the Guest Access module.
	TypeRoomGuestAccess Type = "m.room.guest_access"

	// Events from the Tag module.
	TypeTag = "m.tag"

	// Events from the Room Upgrade module.
	TypeRoomTombstone Type = "m.room.tombstone"
)

var parser = map[Type]func(RawEvent, json.RawMessage) (Event, error){
	TypeRoomCanonicalAlias: defaultParse(func() Event { return new(RoomCanonicalAliasEvent) }),
	TypeRoomCreate:         defaultParse(func() Event { return new(RoomCreateEvent) }),
	TypeRoomJoinRules:      defaultParse(func() Event { return new(RoomJoinRulesEvent) }),
	TypeRoomMember:         parseRoomMemberEvent,
	TypeRoomPowerLevels:    defaultParse(func() Event { return new(RoomPowerLevelsEvent) }),
	TypeRoomRedaction:      parseRoomRedactionEvent,

	TypeRoomMessage: defaultParse(func() Event { return new(RoomMessageEvent) }),
	TypeRoomName:    defaultParse(func() Event { return new(RoomNameEvent) }),
	TypeRoomTopic:   defaultParse(func() Event { return new(RoomTopicEvent) }),
	TypeRoomAvatar:  defaultParse(func() Event { return new(RoomAvatarEvent) }),
	TypeRoomPinned:  defaultParse(func() Event { return new(RoomPinnedEvent) }),

	TypeDirect: parseDirectEvent,

	TypeCallInvite:     defaultParse(func() Event { return new(CallInviteEvent) }),
	TypeCallCandidates: defaultParse(func() Event { return new(CallCandidatesEvent) }),
	TypeCallAnswer:     defaultParse(func() Event { return new(CallAnswerEvent) }),
	TypeCallHangup:     defaultParse(func() Event { return new(CallHangupEvent) }),

	TypeTyping: parseTypingEvent,

	TypeReceipt: parseReceiptEvent,

	TypePresence: parsePresenceEvent,

	TypeRoomHistoryVisibility: defaultParse(func() Event { return new(RoomHistoryVisibilityEvent) }),

	TypeRoomGuestAccess: defaultParse(func() Event { return new(RoomGuestAccessEvent) }),

	TypeTag: defaultParse(func() Event { return new(TagEvent) }),

	TypeRoomTombstone: defaultParse(func() Event { return new(RoomTombstoneEvent) }),
}
