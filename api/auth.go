package api

import (
	"errors"
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// Errors returned by (*Client).Login.
var (
	// ErrInvalidRequest means that an invalid request has been provided.
	// Examples are when an unsupported login method is provided.
	ErrInvalidRequest = errors.New("invalid request provided")
	// ErrInvalidCreds means that invalid credentials have been provided.
	ErrInvalidCreds = errors.New("invalid credentials provided")
	// ErrUserDeactivated means that the user being logged in to has already been deactivated.
	ErrUserDeactivated = errors.New("user being logged in to has been deactivated")
)

// GetLoginMethods return the login methods supported by the homeserver.
//
// It implements the `GET _matrix/client/r0/login` endpoint.
func (c *Client) GetLoginMethods() ([]matrix.LoginMethod, error) {
	var response struct {
		Flows []struct {
			Type matrix.LoginMethod `json:"type"`
		} `json:"flows"`
	}

	err := c.Request("GET", "_matrix/client/r0/login", &response)
	if err != nil {
		return nil, fmt.Errorf("error getting login methods: %w", err)
	}

	result := make([]matrix.LoginMethod, 0, len(response.Flows))
	for _, v := range response.Flows {
		result = append(result, v.Type)
	}
	return result, nil
}

// LoginArg represents all possible login arguments.
type LoginArg struct {
	Type                     matrix.LoginMethod `json:"type"`
	InitialDeviceDisplayName string             `json:"initial_device_display_name,omitempty"`

	// DeviceID should be provided when resuming a session.
	DeviceID matrix.DeviceID `json:"device_id,omitempty"`

	// Identifier and Password is only required when logging in with password.
	Identifier matrix.Identifier `json:"identifier,omitempty"`
	Password   string            `json:"password,omitempty"`

	// Token is only required when logging in with token.
	Token string `json:"token,omitempty"`
}

// Login logs the client into the homeserver with the provided arguments.
//
// It implements the `POST _matrix/client/r0/login` endpoint.
func (c *Client) Login(arg LoginArg) error {
	var resp struct {
		UserID      matrix.UserID         `json:"user_id"`
		AccessToken string                `json:"access_token"`
		DeviceID    matrix.DeviceID       `json:"device_id"`
		WellKnown   DiscoveryInfoResponse `json:"well_known"`
	}
	err := c.Request("POST", "_matrix/client/r0/login", &resp, httputil.WithBody(arg))
	if err != nil {
		return fmt.Errorf("error logging in: %w", matrix.MapAPIError(
			err, matrix.ErrorMap{
				matrix.CodeUnknown:         ErrInvalidRequest,
				matrix.CodeForbidden:       ErrInvalidCreds,
				matrix.CodeUserDeactivated: ErrUserDeactivated,
			}),
		)
	}

	c.UserID = resp.UserID
	c.AccessToken = resp.AccessToken
	c.DeviceID = resp.DeviceID

	return nil
}

// Logout clears the AccessToken field in the client and attempts to invalidate the
// token on the server-side.
//
// It implements the `POST _matrix/client/r0/logout` endpoint.
func (c *Client) Logout() error {
	err := c.Request("POST", "_matrix/client/r0/logout", nil, httputil.WithToken())
	c.AccessToken = ""
	if err != nil {
		return fmt.Errorf("error logging out: %w", err)
	}
	return nil
}

// LogoutAll clears the AccessToken field in the client and attempts to invalidate all
// tokens on the server-side.
//
// It implements the `POST _matrix/client/r0/logout/all` endpoint.
func (c *Client) LogoutAll() error {
	err := c.Request("POST", "_matrix/client/r0/logout/all", nil, httputil.WithToken())
	c.AccessToken = ""
	if err != nil {
		return fmt.Errorf("error logging out all tokens: %w", err)
	}
	return nil
}
