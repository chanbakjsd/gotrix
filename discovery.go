package gomatrix

import (
	"net/http"
	"net/url"

	"github.com/chanbakjsd/gomatrix/api"
)

// New is a helper function that calls NewWithClient with the default HTTP client.
func New(serverName string) (*Client, error) {
	return NewWithClient(http.DefaultClient, serverName)
}

// NewWithClient creates a client with the provided HTTP client.
//
// It attempts to discover the homeserver using the provided server name.
// This allows the host address to be extracted from the user ID and used to
// construct a client in a way that is spec-compliant.
func NewWithClient(httpClient *http.Client, serverName string) (*Client, error) {
	parsed, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}

	apiClient := &api.Client{}
	apiClient.HomeServer = parsed.Host
	info, err := apiClient.DiscoveryInfo()
	if err != nil {
		return nil, err
	}
	homeServerURL, err := url.Parse(info.HomeServer.BaseURL)
	if err != nil {
		return nil, err
	}
	identityServerURL, err := url.Parse(info.IdentityServer.BaseURL)
	if err != nil {
		return nil, err
	}

	apiClient.HomeServer = homeServerURL.Host
	apiClient.IdentityServer = identityServerURL.Host

	return &Client{
		Client: apiClient,
	}, nil
}