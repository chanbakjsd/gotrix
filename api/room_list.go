package api

import (
	"fmt"
	"strconv"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// RoomVisibility returns the visibility of a room on the server's public room directory.
func (c *Client) RoomVisibility(roomID matrix.RoomID) (RoomVisibility, error) {
	var resp struct {
		Visibility RoomVisibility `json:"visibility"`
	}

	err := c.Request("GET", EndpointDirectoryListRoom(roomID), &resp)
	if err != nil {
		return "", fmt.Errorf("error fetching room visibility: %w", err)
	}
	return resp.Visibility, nil
}

// RoomVisibilitySet sets the visibility of a room on the server's public room directory.
func (c *Client) RoomVisibilitySet(roomID matrix.RoomID, newVisibility RoomVisibility) error {
	req := struct {
		Visibility RoomVisibility `json:"visibility"`
	}{newVisibility}

	err := c.Request(
		"GET", EndpointDirectoryListRoom(roomID), nil,
		httputil.WithToken(), httputil.WithJSONBody(req),
	)
	if err != nil {
		return fmt.Errorf("error setting room visibility: %w", err)
	}
	return nil
}

// PublicRoom is a room advertised by the server.
type PublicRoom struct {
	// Always provided.
	JoinedMemberCount int    `json:"num_joined_members"`
	RoomID            string `json:"room_id"`
	WorldReadable     bool   `json:"world_readable,omitempty"`
	GuestCanJoin      bool   `json:"guest_can_join,omitempty"`

	// Optional values.
	Aliases        []string    `json:"aliases"`
	AvatarURL      *matrix.URL `json:"avatar_url,omitempty"`
	CanonicalAlias *string     `json:"canonical_alias,omitempty"`
	Name           *string     `json:"name,omitempty"`
	Topic          *string     `json:"topic,omitempty"`
}

// PublicRoomsResponse is the response to (*Client).PublicRooms.
type PublicRoomsResponse struct {
	Chunk []PublicRoom `json:"chunk"`

	// Pagination token to fetch previous or next batch.
	PrevBatch *string `json:"prev_batch"`
	NextBatch *string `json:"next_batch"`

	TotalRoomCountEstimate *int `json:"total_room_count_estimate"`
}

// PublicRooms list all the public rooms advertised by a server. All parameters are optional.
// 'since' should be a token returned in PublicRoomsResponse.
// 'server' defaults to the homeserver if it is an empty string.
func (c *Client) PublicRooms(limit int, since string, server string) (PublicRoomsResponse, error) {
	req := map[string]string{}
	if limit != 0 {
		req["limit"] = strconv.Itoa(limit)
	}
	if since != "" {
		req["since"] = since
	}
	if server != "" {
		req["server"] = server
	}

	var resp PublicRoomsResponse
	err := c.Request(
		"GET", EndpointPublicRooms, &resp,
		httputil.WithQuery(req),
	)
	if err != nil {
		return PublicRoomsResponse{}, fmt.Errorf("error fetching public rooms: %w", err)
	}
	return resp, nil
}

// PublicRoomsSearchFilter is the filter applied to PublicRoomSearch.
type PublicRoomsSearchFilter struct {
	Keyword *string `json:"generic_search_term,omitempty"`
}

// PublicRoomsSearchArg is the argument to (*Client).PublicRoomSearch.
type PublicRoomsSearchArg struct {
	Limit                int                      `json:"-"`
	Since                string                   `json:"-"`
	Server               string                   `json:"-"`
	Filter               *PublicRoomsSearchFilter `json:"filter,omitempty"`
	IncludeAllNetworks   bool                     `json:"include_all_networks,omitempty"`
	ThirdPartyInstanceID string                   `json:"third_party_instance_id"`
}

// PublicRoomsSearch is equivalent to PublicRooms except it allows the user to provide filters to narrow search result.
func (c *Client) PublicRoomsSearch(arg PublicRoomsSearchArg) (PublicRoomsResponse, error) {
	req := map[string]string{}
	if arg.Limit != 0 {
		req["limit"] = strconv.Itoa(arg.Limit)
	}
	if arg.Since != "" {
		req["since"] = arg.Since
	}
	if arg.Server != "" {
		req["server"] = arg.Server
	}

	var resp PublicRoomsResponse
	err := c.Request(
		"POST", EndpointPublicRooms, &resp,
		httputil.WithQuery(req), httputil.WithJSONBody(arg),
	)
	if err != nil {
		return PublicRoomsResponse{}, fmt.Errorf("error fetching public rooms: %w", err)
	}
	return resp, nil
}
