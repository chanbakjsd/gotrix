package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
)

// TurnServersResponse represents the response to (*Client).TurnServers().
type TurnServersResponse struct {
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	URI        []string `json:"uris"`
	TimeToLive int      `json:"ttl"`
}

// TurnServers returns the list of TURN servers which clients can use to contact the remote party.
// It may error if the homeserver doesn't support the VoIP module or if the request failed.
//
// It implements the `GET _matrix/client/r0/voip/turnServer` endpoint.
func (c *Client) TurnServers() (*TurnServersResponse, error) {
	resp := &TurnServersResponse{}
	err := c.Request("GET", "_matrix/client/r0/voip/turnServer", resp, httputil.WithToken())
	if err != nil {
		return nil, fmt.Errorf("error fetching TURN servers: %w", err)
	}
	return resp, nil
}
