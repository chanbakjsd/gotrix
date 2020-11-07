package httputil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/chanbakjsd/gomatrix/debug"
	"github.com/chanbakjsd/gomatrix/matrix"
)

// Client is a HTTP client that is Matrix API-aware.
type Client struct {
	ClientDriver

	// AccessToken is the token to attach to the request.
	AccessToken string
	// HomeServer is the host part of the homeserver and is treated as the host
	// for all requests.
	HomeServer string
	// HomeServerScheme is the scheme to talk to homeserver on.
	// It is https most of the time.
	HomeServerScheme string
}

// NewClient creates a new Client that uses the default HTTP client.
func NewClient() Client {
	return Client{
		ClientDriver: http.DefaultClient,
	}
}

// NewCustomClient creates a new Client that uses the provided ClientDriver.
func NewCustomClient(d ClientDriver) Client {
	return Client{
		ClientDriver: d,
	}
}

// Request makes the request and returns the result.
//
// It may return any HTTP request errors or a matrix.HTTPError which may possibly
// wrap a matrix.APIError.
func (c *Client) Request(method, route string, to interface{}, mods ...Modifier) error {
	// Generate the request.
	req, err := http.NewRequest(method, c.HomeServerScheme+"://"+c.HomeServer+"/"+route, nil)
	if err != nil {
		return err
	}

	// Apply all the request modifiers.
	for _, v := range mods {
		v(c, req)
	}

	if debug.TraceEnabled {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			panic(err)
		}
		debug.Trace("<<<<\n", string(b))
	}

	// Make the request.
	resp, err := c.Do(req)
	if err != nil {
		debug.Trace(">>>> N/A. Error:", err)
		return err
	}

	if debug.TraceEnabled {
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			panic(err)
		}
		debug.Trace(">>>>\n", string(b))
	}

	defer func() {
		// We honestly don't care about errors closing the body.
		_ = resp.Body.Close()
	}()

	// HTTP OK. Just return the object.
	if resp.StatusCode == http.StatusOK {
		if to == nil {
			return nil
		}
		return json.NewDecoder(resp.Body).Decode(to)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Try to decode into target just in case it is expecting the error message.
	_ = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(to)

	// Try parsing it as an API error.
	var apiError matrix.APIError

	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&apiError)
	if err != nil {
		return matrix.NewHTTPError(resp.StatusCode, err)
	}

	// If it's a rate-limit, we intercept it and retry after the recommended time.
	if apiError.Code == matrix.CodeLimitExceeded {
		debug.Fields(map[string]interface{}{
			"sleep_ms": apiError.RetryAfterMillisecond,
		}).Debug("Being rate-limited by homeserver.")
		time.Sleep(time.Duration(apiError.RetryAfterMillisecond) * time.Millisecond)
		return c.Request(method, route, to, mods...)
	}

	return matrix.NewHTTPError(resp.StatusCode, apiError)
}
