package event

var _ StateEvent = &RoomHistoryVisibilityEvent{}

// HistoryVisibility specifies the group that can view the room history.
type HistoryVisibility string

// A list of possible visibility values. The default is "shared".
const (
	// VisibilityInvited allows members to see history from the moment they were invited until
	// they are no longer invited or in the room.
	VisibilityInvited HistoryVisibility = "invited"

	// VisibilityJoined allows members to see history from the moment they join the room until
	// they are no longer in the room.
	VisibilityJoined HistoryVisibility = "joined"

	// VisibilityShared allows everyone to see all history including users who are not from the room
	// as long as they are a member at some point.
	VisibilityShared HistoryVisibility = "shared"

	// VisibilityWorldReadable allows everyone to see all history including users who were never in the room.
	VisibilityWorldReadable HistoryVisibility = "world_readable"
)

// RoomHistoryVisibilityEvent is an event where the visibility of history is changed.
type RoomHistoryVisibilityEvent struct {
	StateEventInfo `json:"-"`
	Visibility     HistoryVisibility `json:"history_visibility,omitempty"`
}
