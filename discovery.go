package gotrix

import (
	"net/url"

	"github.com/chanbakjsd/gotrix/api/httputil"
)

// Discover is a helper function that calls DiscoverWithClient with the default HTTP client.
func Discover(serverName string) (*Client, error) {
	return NewWithClient(httputil.NewClient(), serverName)
}

// DiscoverWithClient  attempts to discover the homeserver using the provided server name.
//
// This allows the host address to be extracted from the user ID and used to discover the
// homeserver host in a way that is spec-compliant.
func DiscoverWithClient(httpClient httputil.Client, serverName string) (*Client, error) {
	apiClient, err := NewWithClient(httpClient, serverName)
	if err != nil {
		return nil, err
	}

	info, err := apiClient.DiscoveryInfo()
	if err != nil {
		return nil, err
	}

	apiClient, err = NewWithClient(httpClient, info.HomeServer.BaseURL)
	if err != nil {
		return nil, err
	}
	identityServerURL, err := url.Parse(info.IdentityServer.BaseURL)
	if err != nil {
		return nil, err
	}
	if identityServerURL.Scheme == "" {
		identityServerURL.Scheme = "https"
	}

	apiClient.IdentityServer = identityServerURL.Host

	return apiClient, nil
}
