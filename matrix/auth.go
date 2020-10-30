package matrix

import (
	"fmt"
)

// LoginMethod represents a possible login method that can be used to authenticate.
type LoginMethod string

// List of official authentications.
// It can be found here: https://matrix.org/docs/spec/client_server/r0.6.1#authentication-types.
const (
	LoginPassword  LoginMethod = "m.login.password"
	LoginRecaptcha LoginMethod = "m.login.recaptcha"
	LoginToken     LoginMethod = "m.login.token"
	LoginOAuth2    LoginMethod = "m.login.oauth2"
	LoginSSO       LoginMethod = "m.login.sso"
	LoginEmail     LoginMethod = "m.login.email.identity"
	LoginPhone     LoginMethod = "m.login.msisdn"
	LoginDummy     LoginMethod = "m.login.dummy"
)

// List of famous custom authentications.
const (
	LoginAppservice LoginMethod = "uk.half-shot.msc2778.login.application_service"
)

// FallbackURL generates the URL that the application can open in order to finish the auth process.
// This can be used when the auth method is not natively supported by the client.
func (l LoginMethod) FallbackURL(authServerHost string, sessionID string) string {
	return fmt.Sprintf(
		"https://%s/_matrix/client/r0/auth/%s/fallback/web?session=%s",
		authServerHost, string(l), sessionID,
	)
}
