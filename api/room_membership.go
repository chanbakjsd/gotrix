package api

import (
	"errors"
	"net/url"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

var (
	// ErrRemoteDoesNotSupportVersion represents an error when the invited user's homeserver does not support the
	// version of the room the user is invited to.
	//
	// It is returned by (*Client).Invite.
	ErrRemoteDoesNotSupportVersion = errors.New("remote does not support room version the user is invited to")

	// ErrInvalidInvite represents an error where the client tried to invite a banned user or lack the permission
	// to invite anyone in the first place.
	//
	// It is returned by (*Client).Invite.
	ErrInvalidInvite = errors.New("an invalid invite has been issued")
)

// ErrUserStillInRoom represents an error where the client requests the homeserver to forget a room that the user is in.
//
// It is returned by (*Client).Forget.
var ErrUserStillInRoom = errors.New("cannot forget room that the user is in")

// ErrNoKickPerm represents an error where the client tried to kick a user without sufficient permission.
//
// It is returned by (*Client).Kick.
var ErrNoKickPerm = errors.New("cannot kick the user from the room due to insufficient permission")

// ErrNoBanPerm represents an error where the client tried to ban a user without sufficient permission.
//
// It is returned by (*Client).Ban.
var ErrNoBanPerm = errors.New("cannot ban the user from the room due to insufficient permission")

// ErrNoUnbanPerm represents an error where the client tried to ban a user without sufficient permission.
//
// It is returned by (*Client).Unban.
var ErrNoUnbanPerm = errors.New("cannot unban the user from the room due to insufficient permission")

// Rooms returns a list of the user's current rooms.
//
// It implements the `GET /_matrix/client/r0/joined_rooms` endpoint.
func (c *Client) Rooms() ([]matrix.RoomID, error) {
	var resp struct {
		JoinedRooms []matrix.RoomID `json:"joined_rooms"`
	}
	err := c.Request(
		"GET", "_matrix/client/r0/joined_rooms", &resp,
		httputil.WithToken(),
	)
	return resp.JoinedRooms, err
}

// Invite invites the requested user to the specified room ID.
//
// It implements the `POST /_matrix/client/r0/rooms/{roomId}/invite` endpoint.
func (c *Client) Invite(roomID matrix.RoomID, userID matrix.UserID) error {
	body := struct {
		UserID string `json:"user_id"`
	}{string(userID)}
	err := c.Request(
		"POST", "_matrix/client/r0/rooms/"+url.PathEscape(string(roomID))+"/invite", nil,
		httputil.WithToken(), httputil.WithBody(body),
	)
	return matrix.MapAPIError(err, matrix.ErrorMap{
		// matrix.CodeBadJSON, matrix.CodeNotJSON shouldn't happen.
		matrix.CodeUnsupportedRoomVersion: ErrRemoteDoesNotSupportVersion,
		matrix.CodeForbidden:              ErrInvalidInvite,
	})
}

// RoomJoin joins the specified room ID.
//
// It implements the `POST /_matrix/client/r0/rooms/{roomId}/join` endpoint.
func (c *Client) RoomJoin(roomID matrix.RoomID) error {
	err := c.Request(
		"POST", "_matrix/client/r0/rooms/"+url.PathEscape(string(roomID))+"/join", nil,
		httputil.WithToken(),
	)
	return matrix.MapAPIError(err, matrix.ErrorMap{
		matrix.CodeForbidden: ErrRoomNotFound,
	})
}

// TODO: Implement third party invite version of (*Client).RoomJoin.

// RoomLeave leaves the specified room ID.
//
// It implements the `POST /_matrix/client/r0/rooms/{roomId}/leave` endpoint.
func (c *Client) RoomLeave(roomID matrix.RoomID) error {
	err := c.Request(
		"POST", "_matrix/client/r0/rooms/"+url.PathEscape(string(roomID))+"/leave", nil,
		httputil.WithToken(),
	)
	return err
}

// RoomForget tells the homeserver that the user no longer intend to fetch events from the provided
// room. This allows the homeserver to delete the room if every previous member forgets it.
// The client must not be in the room when RoomForget is called.
//
// It implements the `POST /_matrix/client/r0/rooms/{roomId}/forget` endpoint.
func (c *Client) RoomForget(roomID matrix.RoomID) error {
	err := c.Request(
		"POST", "_matrix/client/r0/rooms/"+url.PathEscape(string(roomID))+"/forget", nil,
		httputil.WithToken(),
	)
	if err == nil {
		return nil
	}
	switch matrix.StatusCode(err) {
	case 400:
		return ErrUserStillInRoom
	default:
		return err
	}
}

// Kick kicks the user from the provided room.
//
// It implements the `POST /_matrix/client/r0/rooms/{roomId}/kick` endpoint.
func (c *Client) Kick(roomID matrix.RoomID, userID matrix.UserID, reason string) error {
	param := struct {
		UserID matrix.UserID `json:"user_id"`
		Reason string        `json:"reason,omitempty"`
	}{userID, reason}

	err := c.Request(
		"POST", "_matrix/client/r0/rooms/"+url.PathEscape(string(roomID))+"/kick", nil,
		httputil.WithToken(), httputil.WithBody(param),
	)
	if err == nil {
		return nil
	}
	switch matrix.StatusCode(err) {
	case 403:
		return ErrNoKickPerm
	default:
		return err
	}
}

// Ban bans the user from the provided room.
//
// It implements the `POST /_matrix/client/r0/rooms/{roomId}/ban` endpoint.
func (c *Client) Ban(roomID matrix.RoomID, userID matrix.UserID, reason string) error {
	param := struct {
		UserID matrix.UserID `json:"user_id"`
		Reason string        `json:"reason,omitempty"`
	}{userID, reason}

	err := c.Request(
		"POST", "_matrix/client/r0/rooms/"+url.PathEscape(string(roomID))+"/ban", nil,
		httputil.WithToken(), httputil.WithBody(param),
	)
	if err == nil {
		return nil
	}
	switch matrix.StatusCode(err) {
	case 403:
		return ErrNoBanPerm
	default:
		return err
	}
}

// Unban unbans the user from the provided room.
//
// It implements the `POST /_matrix/client/r0/rooms/{roomId}/unban` endpoint.
func (c *Client) Unban(roomID matrix.RoomID, userID matrix.UserID) error {
	param := struct {
		UserID matrix.UserID `json:"user_id"`
	}{userID}

	err := c.Request(
		"POST", "_matrix/client/r0/rooms/"+url.PathEscape(string(roomID))+"/unban", nil,
		httputil.WithToken(), httputil.WithBody(param),
	)
	if err == nil {
		return nil
	}
	switch matrix.StatusCode(err) {
	case 403:
		return ErrNoUnbanPerm
	default:
		return err
	}
}
