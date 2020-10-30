package api

import (
	"encoding/json"
	"errors"

	"github.com/chanbakjsd/gomatrix/matrix"
)

// ErrInteractiveAuthIncomplete is returned when the result is requested before interactive auth
// is complete.
var ErrInteractiveAuthIncomplete = errors.New("interactive auth has not been completed yet")

// UserInteractiveAuthAPI represents the state that needs to be kept for
// the user interactive authentication API.
//
// It implements https://matrix.org/docs/spec/client_server/r0.6.1#user-interactive-authentication-api.
type UserInteractiveAuthAPI struct {
	// Flows are list of stages that the server requires before allowing.
	Flows []struct {
		Stages []matrix.LoginMethod `json:"stages"`
	} `json:"flows"`

	// Params are list of parameters which contain arbitrary data required to
	// finish the auth flow.
	// Examples are the public key for Captcha.
	Params map[string]json.RawMessage `json:"params"`

	// Session represents the current session of auth that allows the server to
	// keep track of the auth flow.
	Session string `json:"session"`

	// Completed lists all auth successes that the server acknowledges.
	Completed []matrix.LoginMethod `json:"completed"`

	// Error and ErrorCode represents the error encountered which probably means
	// incorrect credentials and similar.
	Error     string           `json:"error"`
	ErrorCode matrix.ErrorCode `json:"errcode"`

	// Request is the function to call to make a request.
	Request         func(req, to interface{}) error `json:"-"`
	SuccessCallback func(json.RawMessage) error     `json:"-"`

	// Result is the result after everything succeeds.
	result *json.RawMessage
}

// Auth attempts to authenticate using the provided information in an attempt to progress in the authentication.
func (u *UserInteractiveAuthAPI) Auth(req interface{}) error {
	var rawMsg json.RawMessage
	err := u.Request(req, rawMsg)
	return u.processResponse(rawMsg, err)
}

// processResponse updates the object's state using the provided raw message.
func (u *UserInteractiveAuthAPI) processResponse(rawMsg json.RawMessage, reqError error) error {
	// If there isn't any errors at all (not 401). The request was successful.
	if reqError == nil {
		u.result = &rawMsg
		if u.SuccessCallback != nil {
			return u.SuccessCallback(rawMsg)
		}
		return nil
	}

	// If there's an error in request and it's not unauthorized (the server requesting to continue auth),
	// we can't handle it.
	if matrix.StatusCode(reqError) != 401 {
		return reqError
	}

	var resp *UserInteractiveAuthAPI
	err := json.Unmarshal(rawMsg, resp)
	if err != nil {
		return err
	}

	*u = *resp
	return nil
}

// IsComplete returns true if the interactive auth is complete.
func (u *UserInteractiveAuthAPI) IsComplete() bool {
	return u.result != nil
}

// Result returns a RawMessage or an error if interactive auth is not complete yet.
func (u *UserInteractiveAuthAPI) Result() (*json.RawMessage, error) {
	if u.IsComplete() {
		return nil, ErrInteractiveAuthIncomplete
	}
	return u.result, nil
}

// authRequest is an aggregate of all supported fields.
// It is meant to be internal. To implement custom login types, use your own struct.
type authRequest struct {
	Type    matrix.LoginMethod `json:"type"`
	Session string             `json:"session"`

	Identifier matrix.Identifier `json:"identifier,omitempty"`
	Password   string            `json:"password,omitempty"`

	Response string `json:"response,omitempty"`

	Token         string `json:"token,omitempty"`
	TransactionID string `json:"txn_id,omitempty"`

	ThreePIDCreds ThreePIDCreds `json:"threepidCreds,omitempty"`
}

// ThreePIDCreds represents three PID information returned from the auth server.
type ThreePIDCreds struct {
	IdentitySessionID   string `json:"sid"`
	ClientSecret        string `json:"client_secret"`
	IdentityServerURL   string `json:"id_server"`
	IdentityAccessToken string `json:"id_access_token"`
}

// AuthPassword is a helper method to login with password.
func (u *UserInteractiveAuthAPI) AuthPassword(id matrix.Identifier, password string) error {
	return u.Auth(authRequest{
		Type:       matrix.LoginPassword,
		Session:    u.Session,
		Identifier: id,
		Password:   password,
	})
}

// AuthRecaptcha is a helper method to login with recaptcha.
func (u *UserInteractiveAuthAPI) AuthRecaptcha(response string) error {
	return u.Auth(authRequest{
		Type:     matrix.LoginRecaptcha,
		Session:  u.Session,
		Response: response,
	})
}

// AuthToken is a helper method to login with token.
// Transaction ID should be a randomly generated ID that is persistent for each request.
func (u *UserInteractiveAuthAPI) AuthToken(token, transactionID string) error {
	return u.Auth(authRequest{
		Type:          matrix.LoginToken,
		Session:       u.Session,
		Token:         token,
		TransactionID: transactionID,
	})
}

// AuthEmail is a helper method to login with email verification.
func (u *UserInteractiveAuthAPI) AuthEmail(threePidCreds ThreePIDCreds) error {
	return u.Auth(authRequest{
		Type:          matrix.LoginEmail,
		Session:       u.Session,
		ThreePIDCreds: threePidCreds,
	})
}

// AuthPhone is a helper method to login with phone verification.
func (u *UserInteractiveAuthAPI) AuthPhone(threePidCreds ThreePIDCreds) error {
	return u.Auth(authRequest{
		Type:          matrix.LoginPhone,
		Session:       u.Session,
		ThreePIDCreds: threePidCreds,
	})
}

// AuthDummy is a helper method to login with the dummy method.
// It does not need any credentials and serve only to separate different flows.
func (u *UserInteractiveAuthAPI) AuthDummy() error {
	return u.Auth(authRequest{
		Type:    matrix.LoginDummy,
		Session: u.Session,
	})
}
