package gomatrix

import (
	"net/http"
	"net/url"

	"github.com/chanbakjsd/gomatrix/api"
)

func New(serverName string) (*Client, error) {
	return NewWithClient(http.DefaultClient, serverName)
}

func NewWithClient(httpClient *http.Client, serverName string) (*Client, error) {
	parsed, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}

	apiClient := &api.Client{
		HomeServer: parsed.Host,
	}
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
