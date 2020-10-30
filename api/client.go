package api

import (
	"errors"

	"github.com/chanbakjsd/gomatrix/api/httputil"
	"github.com/chanbakjsd/gomatrix/matrix"
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
	UserID         string
	DeviceID       string
}

// Whoami queries the homeserver to check if the token is still valid.
// The user ID is returned if it's successful.
//
// This implements the `GET _matrix/client/r0/account/whoami` endpoint.
func (c *Client) Whoami() (string, error) {
	var resp struct {
		UserID string `json:"user_id"`
	}

	err := c.Request(
		"GET", "_matrix/client/r0/account/whoami", &resp,
		httputil.WithQuery(map[string]string{
			"user_id": c.UserID,
		}),
		httputil.WithToken(),
	)

	return resp.UserID, matrix.MapAPIError(err, matrix.ErrorMap{
		matrix.CodeUnknownToken: ErrInvalidToken,
		matrix.CodeForbidden:    ErrTokenAndUserMismatch,
	})
}
