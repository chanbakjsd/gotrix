package api

import (
	"context"
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// Client represents a session that can be established to the server.
// It contains every info the server expects to be persisted on client-side.
type Client struct {
	httputil.Client
	IdentityServer string
	UserID         matrix.UserID
	DeviceID       matrix.DeviceID
}

// WithContext creates a copy of Client that uses the provided context when creating a HTTP request.
func (c Client) WithContext(ctx context.Context) *Client {
	c.Client = c.Client.WithContext(ctx)
	return &c
}

// Whoami queries the homeserver to check if the token is still valid.
// The user ID is returned if it's successful.
func (c *Client) Whoami() (matrix.UserID, error) {
	var resp struct {
		UserID matrix.UserID `json:"user_id"`
	}

	err := c.Request(
		"GET", EndpointAccountWhoami, &resp,
		httputil.WithToken(), httputil.WithQuery(map[string]string{
			"user_id": string(c.UserID),
		}),
	)
	if err != nil {
		return "", fmt.Errorf("error fetching whoami: %w", err)
	}
	return resp.UserID, nil
}

// ServerCapabilities retrieves the homeserver's capabilities.
func (c *Client) ServerCapabilities() (*matrix.Capabilities, error) {
	var resp struct {
		Capabilities *matrix.Capabilities `json:"capabilities"`
	}

	err := c.Request(
		"GET", EndpointCapabilities, &resp,
		httputil.WithToken(),
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching server capabilities: %w", err)
	}

	return resp.Capabilities, nil
}
