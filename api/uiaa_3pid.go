package api

import (
	"errors"
	"fmt"
)

// ErrUnsupportedAuthType is returned when a 3PID auth is attempted on an endpoint that doesn't support it.
var ErrUnsupportedAuthType = errors.New("3PID tokens are unsupported by the current interactive auth session")

// RequestEmailTokenArg represents all possible argument to RequestEmailToken.
type RequestEmailTokenArg struct {
	ClientSecret  string `json:"client_secret"`
	Email         string `json:"email"`
	SendAttempt   int    `json:"send_attempt"`
	NextLink      string `json:"next_link,omitempty"`
	IDServer      string `json:"id_server,omitempty"`
	IDAccessToken string `json:"id_access_token,omitempty"`
}

// RequestEmailTokenResponse represents the response to (*UserInteractiveAuthAPI).RequestEmailToken.
type RequestEmailTokenResponse struct {
	SessionID string `json:"sid"`
	SubmitURL string `json:"submit_url"`
}

// RequestEmailToken requests an email token to be mailed to the specified email address
// for registration purposes.
// It returns the session ID needed in 3PID auth and the submit URL (if applicable).
func (u *UserInteractiveAuthAPI) RequestEmailToken(req RequestEmailTokenArg) (*RequestEmailTokenResponse, error) {
	if u.RequestThreePID == nil {
		return nil, ErrUnsupportedAuthType
	}
	response := &RequestEmailTokenResponse{}
	err := u.RequestThreePID("email", req, response)
	if err != nil {
		return nil, fmt.Errorf("uiaa: error requesting email token: %w", err)
	}
	return response, nil
}

// RequestPhoneTokenArg represents all possible argument to RequestPhoneToken.
type RequestPhoneTokenArg struct {
	ClientSecret  string `json:"client_secret"`
	PhoneNumber   string `json:"phone_number"`
	SendAttempt   int    `json:"send_attempt"`
	NextLink      string `json:"next_link,omitempty"`
	IDServer      string `json:"id_server,omitempty"`
	IDAccessToken string `json:"id_access_token,omitempty"`
}

// RequestPhoneTokenResponse represents the response to (*UserInteractiveAuthAPI).RequestPhoneToken.
type RequestPhoneTokenResponse struct {
	SessionID string `json:"sid"`
	SubmitURL string `json:"submit_url"`
}

// RequestPhoneToken requests an email token to be mailed to the specified email address
// for registration purposes.
// It returns the session ID needed in 3PID auth and the submit URL (if applicable).
func (u *UserInteractiveAuthAPI) RequestPhoneToken(req RequestPhoneTokenArg) (*RequestPhoneTokenResponse, error) {
	if u.RequestThreePID == nil {
		return nil, ErrUnsupportedAuthType
	}
	response := &RequestPhoneTokenResponse{}
	err := u.RequestThreePID("phone", req, response)
	if err != nil {
		return nil, fmt.Errorf("uiaa: error requesting phone token: %w", err)
	}
	return response, nil
}
