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
