package api

import (
	"encoding/json"
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// PasswordChange sends a request to the homeserver to change the password.
// All devices except the current one will be logged out if logoutDevices is set to true.
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
			"POST", EndpointAccountPassword, to,
			httputil.WithToken(), httputil.WithJSONBody(req),
		)
		if err != nil {
			return fmt.Errorf("error changing password: %w", err)
		}
		return nil
	}
	uiaa.RequestThreePID = func(authType string, auth, to interface{}) error {
		return c.Request(
			"POST", EndpointAccountPasswordRequestToken(authType), nil,
			httputil.WithJSONBody(auth),
		)
	}
	err := uiaa.Auth(nil)
	return uiaa, err
}

// DeactivateAccount deactivates the account of the current user.
//
// idServer is the identity server to unbind all of the user's 3PID from.
// It is optional and if not provided, the homeserver is responsible for determining
// the unbind source.
//
// It returns an InteractiveDeactivate object which should be used to interactively
// fulfill authentication requirements of the server.
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
		err := c.Request(
			"POST", EndpointAccountDeactivate, to,
			httputil.WithToken(), httputil.WithJSONBody(req),
		)
		if err != nil {
			return fmt.Errorf("error deactivating account: %w", err)
		}
		return nil
	}
	err := uiaa.Auth(nil)
	return uiaa, err
}

// InteractiveDeactivate is a struct that adds response parsing helper functions onto UserInteractiveAuthAPI.
// To see functions on authenticating, refer to it instead.
type InteractiveDeactivate struct {
	*UserInteractiveAuthAPI
}

// DeactivateResponse formats the Result() as the response of the deactivation.
// It returns an error if there isn't any result yet.
func (i InteractiveDeactivate) DeactivateResponse() (matrix.IDServerUnbindResult, error) {
	msg, err := i.Result()
	if err != nil {
		return "", err
	}

	var result struct {
		IDServerUnbindResult matrix.IDServerUnbindResult `json:"id_server_unbind_result"`
	}
	err = json.Unmarshal(*msg, &result)
	return result.IDServerUnbindResult, err
}
