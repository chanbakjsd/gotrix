package event

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

var parser = map[Type]func(RawEvent) (Event, error){
	TypeRoomCanonicalAlias: roomEventParse(func() eventWithRoomEventInfo { return new(RoomCanonicalAliasEvent) }),
	TypeRoomCreate:         roomEventParse(func() eventWithRoomEventInfo { return new(RoomCreateEvent) }),
	TypeRoomJoinRules:      roomEventParse(func() eventWithRoomEventInfo { return new(RoomJoinRulesEvent) }),
	TypeRoomMember:         parseRoomMemberEvent,
	TypeRoomPowerLevels:    roomEventParse(func() eventWithRoomEventInfo { return new(RoomPowerLevelsEvent) }),
	TypeRoomRedaction:      roomEventParse(func() eventWithRoomEventInfo { return new(RoomRedactionEvent) }),

	TypeRoomMessage: roomEventParse(func() eventWithRoomEventInfo { return new(RoomMessageEvent) }),
	TypeRoomName:    roomEventParse(func() eventWithRoomEventInfo { return new(RoomNameEvent) }),
	TypeRoomTopic:   roomEventParse(func() eventWithRoomEventInfo { return new(RoomTopicEvent) }),
	TypeRoomAvatar:  roomEventParse(func() eventWithRoomEventInfo { return new(RoomAvatarEvent) }),
	TypeRoomPinned:  roomEventParse(func() eventWithRoomEventInfo { return new(RoomPinnedEvent) }),

	TypeDirect: eventParse(func() Event { return new(DirectEvent) }),

	TypeCallInvite:     roomEventParse(func() eventWithRoomEventInfo { return new(CallInviteEvent) }),
	TypeCallCandidates: roomEventParse(func() eventWithRoomEventInfo { return new(CallCandidatesEvent) }),
	TypeCallAnswer:     roomEventParse(func() eventWithRoomEventInfo { return new(CallAnswerEvent) }),
	TypeCallHangup:     roomEventParse(func() eventWithRoomEventInfo { return new(CallHangupEvent) }),

	TypeTyping: parseTypingEvent,

	TypeReceipt: parseReceiptEvent,

	TypePresence: parsePresenceEvent,

	TypeRoomHistoryVisibility: parseHistoryVisibilityEvent,

	TypeRoomGuestAccess: roomEventParse(func() eventWithRoomEventInfo { return new(RoomGuestAccessEvent) }),

	TypeTag: parseTagEvent,

	TypeRoomTombstone: roomEventParse(func() eventWithRoomEventInfo { return new(RoomTombstoneEvent) }),
}
