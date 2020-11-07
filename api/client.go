package api

import (
	"errors"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// Errors returned by (*Client).Whoami.
var (
	ErrInvalidToken         = errors.New("the access token is not recognized")
	ErrTokenAndUserMismatch = errors.New("the stored access token cannot masquerade as the stored user")
)

// Client represents a session that can be established to the server.
// It contains every info the server expects to be persisted on client-side.
type Client struct {
	httputil.Client
	IdentityServer string
	UserID         matrix.UserID
	DeviceID       matrix.DeviceID
}

// Whoami queries the homeserver to check if the token is still valid.
// The user ID is returned if it's successful.
//
// This implements the `GET _matrix/client/r0/account/whoami` endpoint.
func (c *Client) Whoami() (matrix.UserID, error) {
	var resp struct {
		UserID matrix.UserID `json:"user_id"`
	}

	err := c.Request(
		"GET", "_matrix/client/r0/account/whoami", &resp,
		httputil.WithQuery(map[string]string{
			"user_id": string(c.UserID),
		}),
		httputil.WithToken(),
	)

	return resp.UserID, matrix.MapAPIError(err, matrix.ErrorMap{
		matrix.CodeUnknownToken: ErrInvalidToken,
		matrix.CodeForbidden:    ErrTokenAndUserMismatch,
	})
}

// ServerCapabilities retrieves the homeserver's capabilities.
func (c *Client) ServerCapabilities() (matrix.Capabilities, error) {
	var resp struct {
		Capabilities matrix.Capabilities `json:"capabilities"`
	}

	err := c.Request(
		"GET", "_matrix/client/r0/capabilities", &resp,
		httputil.WithToken(),
	)

	return resp.Capabilities, err
}
