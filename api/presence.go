package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// Presence is the presence information of a user as returned by (*Client).Presence.
type Presence struct {
	Presence        matrix.Presence `json:"presence"`
	LastActiveAgo   *int            `json:"last_active_ago,omitempty"`
	StatusMsg       *string         `json:"status_msg,omitempty"`
	CurrentlyActive *bool           `json:"currently_active,omitempty"`
}

// Presence fetches the presence of the requested user.
func (c *Client) Presence(userID matrix.UserID) (Presence, error) {
	var resp Presence
	err := c.Request(
		"GET", EndpointPresenceStatus(userID), &resp,
		httputil.WithToken(),
	)
	if err != nil {
		return Presence{}, fmt.Errorf("error fetching presence: %w", err)
	}
	return resp, nil
}

// PresenceSet sets the presence of the provided user to the provided presence and status message.
func (c *Client) PresenceSet(userID matrix.UserID, presence matrix.Presence, statusMsg string) error {
	req := struct {
		Presence  matrix.Presence `json:"presence"`
		StatusMsg string          `json:"status_msg,omitempty"`
	}{
		Presence:  presence,
		StatusMsg: statusMsg,
	}

	err := c.Request(
		"PUT", EndpointPresenceStatus(userID), nil,
		httputil.WithToken(), httputil.WithJSONBody(req),
	)
	if err != nil {
		return fmt.Errorf("error setting presence: %w", err)
	}
	return nil
}
