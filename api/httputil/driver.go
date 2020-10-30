package httputil

import (
	"net/http"
)

// ClientDriver represents a HTTP client that can make requests to the endpoint.
type ClientDriver interface {
	Do(req *http.Request) (*http.Response, error)
}
