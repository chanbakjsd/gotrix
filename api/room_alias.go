package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// RoomAliasResponse represents the response to (*Client).RoomAlias.
type RoomAliasResponse struct {
	RoomID  matrix.RoomID `json:"room_id"`
	Servers []string      `json:"servers"` // A list of servers that are aware of this room alias.
}

// RoomAlias fetches information about a room alias.
func (c *Client) RoomAlias(alias string) (RoomAliasResponse, error) {
	var resp RoomAliasResponse
	err := c.Request(
		"GET", c.Endpoints.DirectoryRoomAlias(alias), &resp,
		httputil.WithToken(),
	)
	if err != nil {
		return RoomAliasResponse{}, fmt.Errorf("error fetching room alias: %w", err)
	}
	return resp, err
}

// RoomAliases fetches all alias of a given room.
func (c *Client) RoomAliases(roomID matrix.RoomID) ([]string, error) {
	var resp struct {
		Aliases []string `json:"aliases"`
	}
	err := c.Request(
		"GET", c.Endpoints.RoomAliases(roomID), &resp,
		httputil.WithToken(),
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching room aliases: %w", err)
	}
	return resp.Aliases, nil
}

// RoomAliasCreate creates a room alias.
func (c *Client) RoomAliasCreate(alias string, roomID matrix.RoomID) error {
	req := struct {
		RoomID matrix.RoomID `json:"room_id"`
	}{
		RoomID: roomID,
	}
	err := c.Request(
		"PUT", c.Endpoints.DirectoryRoomAlias(alias), nil,
		httputil.WithToken(), httputil.WithJSONBody(req),
	)
	if err != nil {
		return fmt.Errorf("error creating room alias: %w", err)
	}
	return nil
}

// RoomAliasDelete deletes a room alias.
func (c *Client) RoomAliasDelete(alias string) error {
	err := c.Request(
		"DELETE", c.Endpoints.DirectoryRoomAlias(alias), nil,
		httputil.WithToken(),
	)
	if err != nil {
		return fmt.Errorf("error deleting room alias: %w", err)
	}
	return nil
}
