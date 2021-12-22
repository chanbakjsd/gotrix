package event

var _ StateEvent = &RoomGuestAccessEvent{}

// GuestAccess is an enum that decides if a guest can join a room.
type GuestAccess string

// The two possible values of GuestAccess.
const (
	GuestAccessCanJoin   GuestAccess = "can_join"
	GuestAccessForbidden GuestAccess = "forbidden"
)

// RoomGuestAccessEvent is an event that controls whether guest users are allowed to join rooms.
// If the event is not present, it's inferred to be forbidden.
type RoomGuestAccessEvent struct {
	StateEventInfo `json:"-"`
	GuestAccess    GuestAccess `json:"guest_access"`
}
