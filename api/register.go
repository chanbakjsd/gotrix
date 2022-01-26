package api

import (
	"encoding/json"
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// RegisterArg represents arguments for the Register function.
type RegisterArg struct {
	Auth                     interface{}     `json:"auth,omitempty"`
	Username                 string          `json:"username"`
	Password                 string          `json:"password"`
	DeviceID                 matrix.DeviceID `json:"device_id,omitempty"`
	InitialDeviceDisplayName string          `json:"initial_device_display_name,omitempty"`
	InhibitLogin             bool            `json:"inhibit_login,omitempty"`
}

// RegisterResponse represents the success response from the register endpoint.
type RegisterResponse struct {
	UserID      matrix.UserID   `json:"user_id"`
	AccessToken string          `json:"access_token"`
	DeviceID    matrix.DeviceID `json:"device_id"`
}

// Register registers an account on the homeserver with the provided arguments.
// Once the authentication is successful, the client is automatically logged in
// if InhibitLogin is set to false in RegisterArg.
//
// It returns an InteractiveRegister object which should be used to interactively
// fulfill authentication requirements of the server.
func (c *Client) Register(kind string, req RegisterArg) (InteractiveRegister, error) {
	ir := InteractiveRegister{
		UserInteractiveAuthAPI: &UserInteractiveAuthAPI{},
	}

	ir.Request = func(auth, to interface{}) error {
		req.Auth = auth
		err := c.Request(
			"POST", c.Endpoints.Register(), to,
			httputil.WithJSONBody(req), httputil.WithQuery(map[string]string{
				"kind": kind,
			}),
		)
		return fmt.Errorf("error while registering: %w", err)
	}

	ir.RequestThreePID = func(authType string, auth, to interface{}) error {
		return c.Request(
			"POST", c.Endpoints.RegisterRequestToken(authType), to,
			httputil.WithJSONBody(auth),
		)
	}

	ir.SuccessCallback = func(json.RawMessage) error {
		resp, err := ir.RegisterResponse()
		if err != nil {
			return fmt.Errorf("error retrieving registration response: %w", err)
		}
		// If inhibitLogin is set, the homeserver probably does not supply us
		// with the info we want.
		if !req.InhibitLogin {
			c.UserID = resp.UserID
			c.AccessToken = resp.AccessToken
			c.DeviceID = resp.DeviceID
		}
		return nil
	}

	err := ir.Auth(nil)

	return ir, err
}

// InteractiveRegister is a struct that adds response parsing helper functions onto UserInteractiveAuthAPI.
// To see functions on authenticating, refer to it instead.
type InteractiveRegister struct {
	*UserInteractiveAuthAPI
}

// RegisterResponse formats the Result() as a RegisterResponse.
// It returns an error if there isn't any result yet or the response is malformed.
func (i InteractiveRegister) RegisterResponse() (*RegisterResponse, error) {
	msg, err := i.Result()
	if err != nil {
		return nil, err
	}

	result := &RegisterResponse{}
	err = json.Unmarshal(*msg, result)
	return result, err
}

// UsernameAvailable returns if the username is reported as available on the homeserver.
//
// Clients should be aware that this might be racey as registration can take place
// between UsernameAvailable() and actual registration.
func (c *Client) UsernameAvailable(username string) (bool, error) {
	err := c.Request(
		"GET", c.Endpoints.RegisterAvailable(), nil,
		httputil.WithQuery(map[string]string{
			"username": username,
		}),
	)
	if err != nil {
		return false, fmt.Errorf("username is not available: %w", err)
	}
	return true, nil
}
