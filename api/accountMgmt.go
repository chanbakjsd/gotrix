package api

import (
	"encoding/json"
	"errors"

	"github.com/chanbakjsd/gomatrix/api/httputil"
	"github.com/chanbakjsd/gomatrix/matrix"
)

// ErrNewPasswordTooWeak means that the homeserver rejected the new password for being too weak.
// It is returned by (*Client).PasswordChange.
var ErrNewPasswordTooWeak = errors.New("new password is too weak")

// PasswordChange sends a request to the homeserver to change the password.
// All devices except the current one will be logged out if logoutDevices is set to true.
//
// This implements the `POST _matrix/client/r0/account/password` endpoint.
func (c *Client) PasswordChange(newPassword string, logoutDevices bool) (*UserInteractiveAuthAPI, error) {
	var req struct {
		Auth          interface{} `json:"auth,omitempty"`
		NewPassword   string      `json:"new_password"`
		LogoutDevices bool        `json:"logout_devices,omitempty"`
	}

	req.NewPassword = newPassword
	req.LogoutDevices = logoutDevices

	uiaa := &UserInteractiveAuthAPI{}
	uiaa.Request = func(auth, to interface{}) error {
		req.Auth = auth
		err := c.Request(
			"POST", "_matrix/client/r0/account/password", to,
			httputil.WithToken(),
			httputil.WithBody(req),
		)
		return matrix.MapAPIError(err, matrix.ErrorMap{
			matrix.CodeWeakPassword: ErrNewPasswordTooWeak,
		})
	}
	uiaa.RequestThreePID = func(authType string, auth, to interface{}) error {
		return c.Request(
			"POST", "_matrix/client/r0/account/password/"+authType+"/requestToken",
			httputil.WithBody(auth),
		)
	}
	err := uiaa.Auth(nil)
	return uiaa, err
}

// DeactivateResponse represents the success response from the deactivate endpoint.
type DeactivateResponse struct {
	IDServerUnbindResult matrix.IDServerUnbindResult `json:"id_server_unbind_result"`
}

// DeactivateAccount deactivates the account of the current user.
//
// idServer is the identity server to unbind all of the user's 3PID from.
// It is optional and if not provided, the homeserver is responsible for determining
// the unbind source.
func (c *Client) DeactivateAccount(idServer string) (InteractiveDeactivate, error) {
	var req struct {
		Auth     interface{} `json:"auth,omitempty"`
		IDServer string      `json:"id_server"`
	}

	req.IDServer = idServer
	uiaa := InteractiveDeactivate{
		UserInteractiveAuthAPI: &UserInteractiveAuthAPI{},
	}
	uiaa.Request = func(auth, to interface{}) error {
		req.Auth = auth
		return c.Request(
			"POST", "_matrix/client/r0/account/deactivate", to,
			httputil.WithToken(),
			httputil.WithBody(req),
		)
	}
	err := uiaa.Auth(nil)
	return uiaa, err
}

// InteractiveDeactivate is a struct that adds helper functions onto UserInteractiveAuthAPI.
// To see functions on authenticating, refer to it instead.
type InteractiveDeactivate struct {
	*UserInteractiveAuthAPI
}

// DeactivateResponse formats the Result() as a DeactivateResponse.
//
// It returns an error if there isn't any result yet.
func (i InteractiveDeactivate) DeactivateResponse() (*DeactivateResponse, error) {
	msg, err := i.Result()
	if err != nil {
		return nil, err
	}

	var result *DeactivateResponse
	err = json.Unmarshal(*msg, result)
	return result, err
}
