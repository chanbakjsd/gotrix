package gotrix

import (
	// Embedding success HTML to be displayed to the user.
	_ "embed"
	"net"
	"net/http"
	"strconv"

	"github.com/chanbakjsd/gotrix/api"
)

//go:embed sso_success.html
var defaultSuccessHTML []byte

// SuccessHTML is the HTML displayed when the user successfully logs in through SSO.
var SuccessHTML = defaultSuccessHTML

// handleSSOResult creates a function that handles the redirect from Matrix's SSO by forwarding it to the provided
// channel.
func handleSSOResult(src string, success chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("loginToken") == "" {
			http.Redirect(w, r, src, http.StatusFound)
			return
		}

		_, _ = w.Write(SuccessHTML)
		success <- r.FormValue("loginToken")
	}
}

// LoginSSO returns a URL that can be used to log the user in through SSO and starts a HTTP server to listen to the
// response from the homeserver.
// It automatically cleans up the HTTP server when the context the client has expires or when the user logs in
// successfully.
//
// The returned function blocks until the login finishes or is canceled and returns an error if the login was
// unsuccessful.
func (c *Client) LoginSSO() (string, func() error, error) {
	// Manually create a listener so we know what the port is.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", nil, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	url := "http://127.0.0.1:" + strconv.Itoa(port) + "/"
	ssoURL := c.FullRoute(api.EndpointSSOLogin(url))

	success := make(chan string)
	errChannel := make(chan error)
	srv := http.Server{
		Handler: handleSSOResult(ssoURL, success),
	}
	go srv.Serve(listener)

	go func() {
		defer srv.Close()

		var token string
		if c.ctx != nil {
			select {
			case <-c.ctx.Done():
				errChannel <- c.ctx.Err()
				return
			case token = <-success:
			}
		} else {
			token = <-success
		}

		errChannel <- c.LoginToken(token)
	}()

	listenResult := func() error {
		return <-errChannel
	}

	return ssoURL, listenResult, nil
}
