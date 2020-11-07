package api

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/chanbakjsd/gotrix/matrix"
)

// Errors returned by (*Client).DiscoveryInfo.
var (
	// ErrServerNotFound represents a IGNORE or FAIL_PROMPT meaning that
	// more information should be requested.
	ErrServerNotFound = errors.New("server not found")
	// ErrDiscoveryFail represents a FAIL_ERROR meaning that the server is
	// Matrix-aware but has returned invalid data.
	// Matrix specs recommend clients to prompt the user for further action
	// in this case.
	ErrDiscoveryFail = errors.New("auto-discovery failed")
)

// DiscoveryInfoResponse represents the response to (*Client).DiscoveryInfo.
type DiscoveryInfoResponse struct {
	HomeServer struct {
		BaseURL string `json:"base_url"`
	} `json:"m.homeserver"`
	IdentityServer struct {
		BaseURL string `json:"base_url"`
	} `json:"m.identity_server"`
}

// DiscoveryInfo discovers homeserver and identity server from the URL set in (*Client).HomeServer
// and validates them.
//
// It implements https://matrix.org/docs/spec/client_server/r0.6.1#well-known-uri.
func (c *Client) DiscoveryInfo() (*DiscoveryInfoResponse, error) {
	// Check well-known URI.
	result := &DiscoveryInfoResponse{}
	err := c.Request("GET", ".well-known/matrix/client", result)
	if err != nil {
		switch matrix.StatusCode(err) {
		case -1:
			return nil, err
		case http.StatusNotFound:
			return nil, ErrServerNotFound
		}
		return nil, ErrDiscoveryFail
	}

	// Check that response is valid.
	if result.HomeServer.BaseURL == "" {
		return nil, ErrServerNotFound
	}
	parsedURL, err := url.Parse(result.HomeServer.BaseURL)
	if err != nil {
		return nil, ErrDiscoveryFail
	}

	// Probe provided homeserver to make sure it's valid.
	checkClient := &Client{}
	checkClient.HomeServer = parsedURL.Host

	_, err = checkClient.SupportedVersions()
	if err != nil {
		return nil, err
	}

	//TODO: Check identity server when it's implemented.

	return result, nil
}

// SupportedVersionsResponse represents the response to (*Client).SupportedVersions.
type SupportedVersionsResponse struct {
	Versions         []string        `json:"versions"`
	UnstableFeatures map[string]bool `json:"unstable_features"`
}

// SupportedVersions returns the list of versions supported by a homeserver.
//
// The homeserver is inferred from (*Client).HomeServer and should be set before calling this function.
//
// It implements the `GET _matrix/client/versions` endpoint.
func (c *Client) SupportedVersions() (*SupportedVersionsResponse, error) {
	result := &SupportedVersionsResponse{}
	err := c.Request("GET", "_matrix/client/versions", &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
