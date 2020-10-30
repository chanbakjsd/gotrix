package api

import (
	"github.com/chanbakjsd/gomatrix/api/httputil"
)

// Client represents a session that can be established to the server.
// It contains every info the server expects to be persisted on client-side.
type Client struct {
	httputil.Client
	IdentityServer string
	UserID         string
	DeviceID       string
}
