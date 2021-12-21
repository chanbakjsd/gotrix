package event

import (
	"github.com/chanbakjsd/gotrix/matrix"
)

var (
	_ StateEvent = RoomCanonicalAliasEvent{}
	_ StateEvent = RoomCreateEvent{}
	_ StateEvent = RoomJoinRulesEvent{}
	_ StateEvent = RoomMemberEvent{}
	_ StateEvent = RoomPowerLevelsEvent{}
	_ RoomEvent  = RoomRedactionEvent{}
)

// RoomCanonicalAliasEvent represents a state event where the alias (name) of the room is set.
type RoomCanonicalAliasEvent struct {
	*StateEventInfo

	// The canonical alias for the room. May be empty.
	Alias string `json:"alias,omitempty"`
	// Alternative aliases the room advertises. It can be present even if Alias is empty.
	AltAlias []string `json:"alt_aliases,omitempty"`
}

// RoomCreateEvent represents a state event where the room is created or upgraded.
// Do note that there's no order of Matrix version and it is still considered upgrading for
// "upgrading" version 2 to 1.
// It is the first event in any room.
type RoomCreateEvent struct {
	*StateEventInfo

	// The user ID of the room creator. This is set by the homeserver.
	Creator matrix.UserID `json:"creator"`
	// Whether users from other servers can join. Defaults to true.
	Federated *bool `json:"m.federate,omitempty"`
	// Room Version. Defaults to "1" if not specified.
	RoomVersion *string `json:"room_version,omitempty"`
}

// RoomJoinRulesEvent represents a state event where the room's join rules are set.
type RoomJoinRulesEvent struct {
	*StateEventInfo

	// The new rules to be applied to users wishing to join the room.
	JoinRule JoinRule `json:"join_rule"`
}

// JoinRule represents the condition required to join a room.
type JoinRule string

// "public" means the room can be joined by everyone while "invite" means the user must be
// invited before attempting to join.
//
// "knock" and "private" are reserved keywords which are not implemented.
const (
	JoinPublic  JoinRule = "public"
	JoinKnock   JoinRule = "knock"
	JoinInvite  JoinRule = "invite"
	JoinPrivate JoinRule = "private"
)

// RoomMemberEvent represents a state event where a user's membership state changes.
type RoomMemberEvent struct {
	*StateEventInfo

	// The ID of the user for this event.
	UserID matrix.UserID `json:"-"`
	// The avatar URL of the user, if any.
	AvatarURL matrix.URL `json:"avatar_url,omitempty"`
	// The display name of the user, if any.
	DisplayName *string `json:"displayname,omitempty"`
	// The new state of the user in the room.
	NewState MemberType `json:"membership,omitempty"`
	// Flag indicating if the room was created with intention of being a DM.
	IsDirect bool `json:"is_direct,omitempty"`
	// ThirdPartyInvites is set when it's an invite event and is the successor of a
	// m.room.third_party_invite event.
	ThirdPartyInvite struct {
		DisplayName string `json:"display_name"`
	} `json:"third_party_invite,omitempty"`
	// A purely INFORMATIONAL source that SHOULD NOT be trusted for the state of the room.
	// It may be present or absent.
	Unsigned struct {
		InviteRoomState []StrippedEvent `json:"invite_room_state"`
	} `json:"unsigned,omitempty"`
}

// MemberType represents the type of member the user is in a room.
type MemberType string

// Invited means that the user is invited and could join the room.
// Joined means that the user is already in the room.
// Left means that the user has not joined the room/left it.
// Banned means that the user has been banned.
//
// Knock is reserved and not implemented.
const (
	MemberInvited MemberType = "invite"
	MemberJoined  MemberType = "join"
	MemberLeft    MemberType = "leave"
	MemberBanned  MemberType = "ban"
	MemberKnock   MemberType = "knock"
)

// RoomPowerLevelsEvent represents a state event that establishes the power level and requirements
// for each event to be sent.
type RoomPowerLevelsEvent struct {
	*StateEventInfo

	// Ban, invite, kick and redact defaults to 50 if unspecified.
	BanRequirement    *int `json:"ban,omitempty"`
	InviteRequirement *int `json:"invite,omitempty"`
	KickRequirement   *int `json:"kick,omitempty"`
	RedactRequirement *int `json:"redact,omitempty"`

	// The power requirements of events. Events overrides the default.
	// The default for normal events is EventRequirement and
	// the default for state events is StateRequirement.
	Events           map[Type]int `json:"events,omitempty"`
	EventRequirement int          `json:"events_default,omitempty"`
	StateRequirement int          `json:"state_default,omitempty"`

	// UserLevel is a map of user IDs to their power level.
	UserLevel map[matrix.UserID]int `json:"users,omitempty"`
	// The default power level of users (if not in UserLevel).
	UserDefault int `json:"users_default,omitempty"`

	Notifications struct {
		// The power level required to ping a room. Defaults to 50.
		Room *int `json:"room,omitempty"`
	} `json:"notifications,omitempty"`
}

// RoomRedactionEvent is a message event where another event is redacted from the history.
// All keys associated with the event may be stripped off, causing the data to no longer be
// accessible.
// This can also be used for moderators to hide message events (which can be undone).
type RoomRedactionEvent struct {
	*RoomEventInfo

	Redacts matrix.EventID `json:"redacts,omitempty"`
	Reason  string         `json:"reason,omitempty"`
}
