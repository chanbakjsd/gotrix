package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/chanbakjsd/gomatrix/matrix"
)

// Client represents a session that can be established to the server.
// It contains every info the server expects to be persisted on client-side.
type Client struct {
	Client         http.Client
	HomeServer     string
	IdentityServer string
	AccessToken    string
	UserID         string
	DeviceID       string
}

// Modifier modifies the Request before it's sent out to add extra info.
type Modifier func(c *Client, req *http.Request)

// ErrorMap represents a map of internal error codes to user-friendly errors.
type ErrorMap map[matrix.ErrorCode]error

// Request processes the method and path to send to and make a request.
// ErrorMap may be provided to convert internal REST errors to user-friendly errors.
// If it's not provided or an unknown error is returned, it'll return the REST error as-is.
// Modifiers may be provided to modify the request before it's sent out.
func (c *Client) Request(method string, route string, to interface{}, errors ErrorMap, mods ...Modifier) error {
	// Generate the request.
	fullRoute := fmt.Sprintf("https://%s/%s", c.HomeServer, route)
	req, err := http.NewRequest(method, fullRoute, nil)
	if err != nil {
		return err
	}

	// Apply all the request modifiers.
	for _, v := range mods {
		v(c, req)
	}

	// Make the request.
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 200 OK. Just return the error.
	if resp.StatusCode == 200 {
		if to == nil {
			return nil
		}
		return json.NewDecoder(resp.Body).Decode(to)
	}

	// Try parsing it as an API error.
	var apiError Error
	apiError.ResponseCode = resp.StatusCode

	err = json.NewDecoder(resp.Body).Decode(&apiError)
	if err != nil {
		return HTTPError{
			Code:            resp.StatusCode,
			UnderlyingError: err,
		}
	}

	// If it's a rate-limit, we intercept it and retry after the recommended time.
	if apiError.Code == matrix.CodeLimitExceeded {
		time.Sleep(time.Duration(apiError.RetryAfterMillisecond) * time.Millisecond)
		return c.Request(method, route, to, errors, mods...)
	}

	// Map error codes to error.
	if x, ok := errors[apiError.Code]; ok {
		return x
	}

	return apiError
}

// WithToken attaches the AccessToken to the request.
// It should be included with requests that require authentication.
func WithToken() Modifier {
	return func(c *Client, req *http.Request) {
		req.Header.Add("Authorization", "Bearer "+c.AccessToken)
	}
}

// WithBody attaches a body to the request. This is generally a byte buffer written to by json.Encoder.
func WithBody(body io.ReadCloser) Modifier {
	return func(_ *Client, req *http.Request) {
		req.Body = body
	}
}

// WithQuery attaches one-to-one queries to the request.
// It is provided as a helper function that calls WithFullQuery.
func WithQuery(rawQueries map[string]string) Modifier {
	fullQuery := make(map[string][]string)
	for k, v := range rawQueries {
		fullQuery[k] = []string{v}
	}
	return WithFullQuery(fullQuery)
}

// WithFullQuery attaches one-to-many queries to the request.
func WithFullQuery(query map[string][]string) Modifier {
	encoded := url.Values(query).Encode()
	return func(_ *Client, req *http.Request) {
		req.URL.RawQuery = encoded
	}
}
