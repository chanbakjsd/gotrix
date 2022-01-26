package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// UpgradeRoom creates a new room linked to the specified room with the specified new version.
func (c *Client) UpgradeRoom(roomID matrix.RoomID, newVersion string) (matrix.RoomID, error) {
	req := map[string]string{
		"new_version": newVersion,
	}
	var resp struct {
		ReplacementRoom matrix.RoomID `json:"replacement_room"`
	}
	err := c.Request(
		"POST", c.Endpoints.RoomUpgrade(roomID), &resp,
		httputil.WithToken(), httputil.WithJSONBody(req),
	)
	if err != nil {
		return "", fmt.Errorf("error upgrading room: %w", err)
	}
	return resp.ReplacementRoom, nil
}
