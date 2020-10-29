package api

import (
	"errors"

	"github.com/chanbakjsd/gomatrix/matrix"
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

	err := c.Request("GET", "_matrix/client/r0/login", &response, nil)
	if err != nil {
		return nil, err
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
	InitialDeviceDisplayName string             `json:"initial_device_display_name"`

	// DeviceID should be provided when resuming a session.
	DeviceID string `json:"device_id"`

	// Identifier and Password is only required when logging in with password.
	Identifier matrix.Identifier `json:"identifier"`
	Password   string            `json:"password"`

	// Token is only required when logging in with token.
	Token string `json:"token"`
}

// Login logs the client into the homeserver with the provided arguments.
//
// It implements the `POST _matrix/client/r0/login` endpoint.
func (c *Client) Login(arg LoginArg) error {
	var resp struct {
		UserID      string                `json:"user_id"`
		AccessToken string                `json:"access_token"`
		DeviceID    string                `json:"device_id"`
		WellKnown   DiscoveryInfoResponse `json:"well_known"`
	}
	err := c.Request(
		"POST", "_matrix/client/r0/login", &resp,
		ErrorMap{
			matrix.CodeUnknown:         ErrInvalidRequest,
			matrix.CodeForbidden:       ErrInvalidCreds,
			matrix.CodeUserDeactivated: ErrUserDeactivated,
		},
	)
	if err != nil {
		return err
	}

	c.UserID = resp.UserID
	c.AccessToken = resp.AccessToken
	c.DeviceID = resp.DeviceID

	return nil
}

// Logout clears the AccessToken field in the client and attempts to invalidate the
// token on the server-side.
//
// It implements the `GET _matrix/client/r0/logout` endpoint.
//
// See (*Client).LogoutAll() for invalidating all active tokens.
func (c *Client) Logout() error {
	err := c.Request("GET", "_matrix/client/r0/logout", nil, nil, WithToken())
	c.AccessToken = ""
	return err
}

// LogoutAll clears the AccessToken field in the client and attempts to invalidate all
// tokens on the server-side.
//
// It implements the `GET _matrix/client/r0/logout/all` endpoint.
//
// See (*Client).Logout() for invalidating only the current token.
func (c *Client) LogoutAll() error {
	err := c.Request("GET", "_matrix/client/r0/logout/all", nil, nil, WithToken())
	c.AccessToken = ""
	return err
}
