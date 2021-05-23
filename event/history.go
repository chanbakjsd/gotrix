package event

import "encoding/json"

var _ StateEvent = RoomHistoryVisibilityEvent{}

// HistoryVisibility specifies the group that can view the room history.
type HistoryVisibility string

const (
	// VisibilityInvited allows members to see history from the moment they were invited until
	// they are no longer invited or in the room.
	VisibilityInvited HistoryVisibility = "invited"
	// VisibilityJoined allows members to see history from the moment they join the room until
	// they are no longer in the room.
	VisibilityJoined HistoryVisibility = "joined"
	// VisibilityShared allows everyone to see all history including users who are not from the room
	// as long as they are a member at some point.
	// ! THIS IS THE DEFAULT !
	VisibilityShared HistoryVisibility = "shared"
	// VisibilityWorldReadable allows everyone to see all history including users who were never in the room.
	VisibilityWorldReadable HistoryVisibility = "world_readable"
)

// RoomHistoryVisibilityEvent is an event where the visibility of history is changed.
type RoomHistoryVisibilityEvent struct {
	RoomEventInfo
	Visibility HistoryVisibility `json:"history_visibility,omitempty"`
}

// Type implements StateEvent.
func (RoomHistoryVisibilityEvent) Type() Type {
	return TypeRoomHistoryVisibility
}

// StateKey implements StateEvent.
func (RoomHistoryVisibilityEvent) StateKey() string {
	return ""
}

func parseHistoryVisibilityEvent(e RawEvent) (Event, error) {
	c := RoomHistoryVisibilityEvent{
		RoomEventInfo: e.toRoomEventInfo(),
	}
	err := json.Unmarshal(e.Content, &c)
	return c, err
}
