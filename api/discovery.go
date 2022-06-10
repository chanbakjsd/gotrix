package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

// DiscoveryInfo discovers homeserver and identity server from the URL set in (*Client).HomeServer.
//
// It implements https://spec.matrix.org/v1.2/client-server-api/#well-known-uri.
func (c *Client) DiscoveryInfo() (*DiscoveryInfoResponse, error) {
	// Check well-known URI.
	var result DiscoveryInfoResponse
	err := c.Request("GET", ".well-known/matrix/client", &result)
	if err != nil {
		switch matrix.StatusCode(err) {
		case -1:
			return nil, fmt.Errorf("error fetching discovery info: %w", err)
		case http.StatusNotFound:
			return nil, ErrServerNotFound
		}
		return nil, ErrDiscoveryFail
	}

	// Check that response is valid.
	if result.HomeServer.BaseURL == "" {
		return nil, ErrServerNotFound
	}

	// TODO: Check identity server when it's implemented.

	return &result, nil
}

// SupportedVersionsResponse represents the response to (*Client).SupportedVersions.
type SupportedVersionsResponse struct {
	Versions         []string        `json:"versions"`
	UnstableFeatures map[string]bool `json:"unstable_features"`
}

// SupportedVersions returns the list of versions supported by a homeserver.
//
// The homeserver is inferred from (*Client).HomeServer and should be set before calling this function.
func (c *Client) SupportedVersions() (SupportedVersionsResponse, error) {
	var result SupportedVersionsResponse
	err := c.Request("GET", EndpointSupportedVersions, &result)
	if err != nil {
		return SupportedVersionsResponse{}, fmt.Errorf("error fetching homeserver supported versions: %w", err)
	}

	return result, nil
}

// WithLatestVersion returns the client that uses the latest endpoint version. An error with a nil
// *Client is returned if the server doesn't have any supported versions or if the server returns
// invalid versions.
func (c Client) WithLatestVersion() (*Client, error) {
	versions, err := c.SupportedVersions()
	if err != nil {
		return nil, err
	}

	var endpointVer string

	for _, v := range versions.Versions {
		base, ok := SupportedVersions[v]
		if !ok {
			continue
		}

		// Always pick the later version if possible.
		if endpointVer == "" || versionLess(endpointVer, base) {
			endpointVer = base
		}
	}

	if endpointVer == "" {
		return nil, fmt.Errorf("server has no supported version in %q", versions.Versions)
	}

	c.Endpoints.Version = endpointVer
	return &c, nil
}

// versionLess returns true if ver1 < ver2 in Matrix endpoint versioning format.
func versionLess(ver1, ver2 string) bool {
	if len(ver1) == 0 || len(ver2) == 0 {
		return ver1 < ver2
	}

	if ver1[0] == ver2[0] {
		v1, err1 := strconv.Atoi(ver1[1:])
		v2, err2 := strconv.Atoi(ver2[1:])
		if err1 == nil && err2 == nil {
			return v1 < v2
		}
		return ver1[1:] < ver2[1:]
	}

	if ver1[0] == 'r' {
		return true // 1 < 2
	}
	if ver2[0] == 'r' {
		return false // 1 > 2
	}

	return ver1 < ver2
}
