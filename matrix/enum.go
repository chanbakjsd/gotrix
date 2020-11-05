package matrix

// IDServerUnbindResult represents whether 3PIDs has been unbound from the
// identity server successfully.
type IDServerUnbindResult string

const (
	// UnbindSuccess represents success in unbinding.
	UnbindSuccess IDServerUnbindResult = "success"
	// UnbindNoSupport means that the homeserver is unable to determine
	// the identity server to unbind from.
	UnbindNoSupport IDServerUnbindResult = "no-support"
)

// Presence represents the status of the client to set while the client is
// polling.
type Presence string

// The three possible status to be in are online, offline (invisible) and
// idle (unavailable).
const (
	PresenceOnline  Presence = "online"
	PresenceOffline Presence = "offline"
	PresenceIdle    Presence = "unavailable"
)
