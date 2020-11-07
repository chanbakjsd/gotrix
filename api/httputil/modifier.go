package httputil

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Modifier modifies the Request before it's sent out to add extra info.
type Modifier func(c *Client, req *http.Request)

// WithToken attaches the AccessToken to the request.
// It should be included with requests that require authentication.
func WithToken() Modifier {
	return func(c *Client, req *http.Request) {
		req.Header.Add("Authorization", "Bearer "+c.AccessToken)
	}
}

// WithBody attaches a JSON body to the request.
func WithBody(body interface{}) Modifier {
	return func(_ *Client, req *http.Request) {
		rp, wp := io.Pipe()
		go func() {
			err := json.NewEncoder(wp).Encode(&body)
			if err != nil {
				panic(err)
			}
			err = wp.Close()
			if err != nil {
				panic(err)
			}
		}()

		req.Header.Add("Content-Type", "application/json")
		req.Body = rp
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
