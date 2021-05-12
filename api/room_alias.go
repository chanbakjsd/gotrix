package api

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// ErrRoomAliasNotFound represents an error where the requested room alias is not found.
// It is returned by (*Client).RoomAlias and (*Client).RoomAliasDelete.
var ErrRoomAliasNotFound = errors.New("requested room alias is not found")

// ErrRoomAliasAlreadyExists represents an error where the room alias requested to be created already exist.
// It is returned by (*Client).RoomAliasCreate.
var ErrRoomAliasAlreadyExists = errors.New("room alias already exists")

// RoomAliasResponse represents the response to (*Client).RoomAlias.
type RoomAliasResponse struct {
	RoomID  matrix.RoomID `json:"room_id"`
	Servers []string      `json:"servers"` // A list of servers that are aware of this room alias.
}

// RoomAlias fetches information about a room alias.
//
// It implements the `GET /_matrix/client/r0/directory/room/{roomAlias}` endpoint.
func (c *Client) RoomAlias(alias string) (*RoomAliasResponse, error) {
	var resp RoomAliasResponse
	err := c.Request(
		"GET", "_matrix/client/r0/directory/room/"+url.PathEscape(alias), &resp,
		httputil.WithToken(),
	)
	if err != nil {
		return nil, matrix.MapAPIError(err, matrix.ErrorMap{
			matrix.CodeNotFound: ErrRoomAliasNotFound,
		})
	}
	return &resp, err
}

// RoomAliases fetches all alias of a given room.
//
// It implements the `GET /_matrix/client/r0/rooms/{roomId}/aliases` endpoint.
func (c *Client) RoomAliases(roomID matrix.RoomID) ([]string, error) {
	var resp struct {
		Aliases []string `json:"aliases"`
	}
	err := c.Request(
		"GET", "_matrix/client/r0/rooms/"+url.PathEscape(string(roomID))+"/aliases", &resp,
		httputil.WithToken(),
	)
	if err != nil {
		return nil, matrix.MapAPIError(err, matrix.ErrorMap{
			matrix.CodeForbidden: ErrRoomNotFound,
		})
	}
	return resp.Aliases, nil
}

// RoomAliasCreate creates a room alias.
//
// It implements the `PUT /_matrix/client/r0/directory/room/{roomAlias}` endpoint.
func (c *Client) RoomAliasCreate(alias string, roomID matrix.RoomID) error {
	req := struct {
		RoomID matrix.RoomID `json:"room_id"`
	}{
		RoomID: roomID,
	}
	err := c.Request(
		"PUT", "_matrix/client/r0/directory/room/"+url.PathEscape(alias), nil,
		httputil.WithToken(), httputil.WithJSONBody(req),
	)

	if matrix.StatusCode(err) == http.StatusGone {
		return ErrRoomAliasAlreadyExists
	}
	return err
}

// RoomAliasDelete deletes a room alias.
//
// It implements the `DELETE /_matrix/client/r0/directory/room/{roomAlias}` endpoint.
func (c *Client) RoomAliasDelete(alias string, roomID matrix.RoomID) error {
	err := c.Request(
		"DELETE", "_matrix/client/r0/directory/room/"+url.PathEscape(alias), nil,
		httputil.WithToken(),
	)

	return matrix.MapAPIError(err, matrix.ErrorMap{
		matrix.CodeNotFound: ErrRoomAliasNotFound,
	})
}
