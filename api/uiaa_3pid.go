package api

import (
	"errors"

	"github.com/chanbakjsd/gomatrix/matrix"
)

// Errors returned by (*UserInteractiveAuthAPI).RequestEmailToken or (*UserInteractiveAuthAPI).RequestPhoneToken.
var (
	ErrEmailAddressInUse = errors.New("requested email address is already in use")
	ErrPhoneNumberInUse  = errors.New("requested phone number is already in use")
	ErrInvalidIDServer   = errors.New("the requested identity server is not trusted by the server")
	ErrThreePIDDisabled  = errors.New("the homeserver does not support third-party identifiers")
)

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
	var response *RequestEmailTokenResponse
	err := u.RequestThreePID("email", req, response)
	return response, matrix.MapAPIError(err, matrix.ErrorMap{
		matrix.CodeThreePIDInUse:    ErrEmailAddressInUse,
		matrix.CodeServerNotTrusted: ErrInvalidIDServer,
		matrix.CodeThreePIDDenied:   ErrThreePIDDisabled,
	})
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
	var response *RequestPhoneTokenResponse
	err := u.RequestThreePID("phone", req, response)
	return response, matrix.MapAPIError(err, matrix.ErrorMap{
		matrix.CodeThreePIDInUse:    ErrPhoneNumberInUse,
		matrix.CodeServerNotTrusted: ErrInvalidIDServer,
		matrix.CodeThreePIDDenied:   ErrThreePIDDisabled,
	})
}
