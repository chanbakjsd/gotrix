package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
)

// TURNServersResponse represents the response to (*Client).TURNServers().
type TURNServersResponse struct {
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	URI        []string `json:"uris"`
	TimeToLive int      `json:"ttl"`
}

// TURNServers returns the list of TURN servers which clients can use to contact the remote party.
// It may error if the homeserver doesn't support the VoIP module or if the request failed.
func (c *Client) TURNServers() (TURNServersResponse, error) {
	var resp TURNServersResponse
	err := c.Request(
		"GET", c.Endpoints.VOIPTURNServers(), resp,
		httputil.WithToken(),
	)
	if err != nil {
		return TURNServersResponse{}, fmt.Errorf("error fetching TURN servers: %w", err)
	}
	return resp, nil
}
