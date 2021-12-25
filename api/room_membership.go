package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// Rooms returns a list of the user's current rooms.
func (c *Client) Rooms() ([]matrix.RoomID, error) {
	var resp struct {
		JoinedRooms []matrix.RoomID `json:"joined_rooms"`
	}
	err := c.Request(
		"GET", EndpointJoinedRooms, &resp,
		httputil.WithToken(),
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching room list: %w", err)
	}
	return resp.JoinedRooms, nil
}

// Invite invites the requested user to the specified room ID.
// To not specify a reason, pass an empty string to the reason parameter.
func (c *Client) Invite(roomID matrix.RoomID, userID matrix.UserID, reason string) error {
	body := struct {
		UserID matrix.UserID `json:"user_id"`
		Reason string        `json:"reason,omitempty"`
	}{userID, reason}
	err := c.Request(
		"POST", EndpointRoomInvite(roomID), nil,
		httputil.WithToken(), httputil.WithJSONBody(body),
	)
	if err != nil {
		return fmt.Errorf("error inviting user into room: %w", err)
	}
	return nil
}

// RoomJoin joins the specified room ID.
// To not specify a reason, pass an empty string to the reason parameter.
func (c *Client) RoomJoin(roomID matrix.RoomID, reason string) error {
	body := struct {
		Reason string `json:"reason,omitempty"`
	}{reason}
	err := c.Request(
		"POST", EndpointRoomJoin(roomID), nil,
		httputil.WithToken(), httputil.WithJSONBody(body),
	)
	if err != nil {
		return fmt.Errorf("error joining room: %w", err)
	}
	return nil
}

// TODO: Implement third party invite version of (*Client).RoomJoin.

// RoomLeave leaves the specified room ID.
// To not specify a reason, pass an empty string to the reason parameter.
func (c *Client) RoomLeave(roomID matrix.RoomID, reason string) error {
	body := struct {
		Reason string `json:"reason,omitempty"`
	}{reason}
	err := c.Request(
		"POST", EndpointRoomLeave(roomID), nil,
		httputil.WithToken(), httputil.WithJSONBody(body),
	)
	return err
}

// RoomForget tells the homeserver that the user no longer intend to fetch events from the provided
// room. This allows the homeserver to delete the room if every previous member forgets it.
// The client must not be in the room when RoomForget is called.
func (c *Client) RoomForget(roomID matrix.RoomID) error {
	err := c.Request(
		"POST", EndpointRoomForget(roomID), nil,
		httputil.WithToken(),
	)
	if err != nil {
		return fmt.Errorf("error forgetting room: %w", err)
	}
	return nil
}

// Kick kicks the user from the provided room.
func (c *Client) Kick(roomID matrix.RoomID, userID matrix.UserID, reason string) error {
	param := struct {
		UserID matrix.UserID `json:"user_id"`
		Reason string        `json:"reason,omitempty"`
	}{userID, reason}

	err := c.Request(
		"POST", EndpointRoomKick(roomID), nil,
		httputil.WithToken(), httputil.WithJSONBody(param),
	)
	if err != nil {
		return fmt.Errorf("error kicking user: %w", err)
	}
	return nil
}

// Ban bans the user from the provided room.
func (c *Client) Ban(roomID matrix.RoomID, userID matrix.UserID, reason string) error {
	param := struct {
		UserID matrix.UserID `json:"user_id"`
		Reason string        `json:"reason,omitempty"`
	}{userID, reason}

	err := c.Request(
		"POST", EndpointRoomBan(roomID), nil,
		httputil.WithToken(), httputil.WithJSONBody(param),
	)
	if err != nil {
		return fmt.Errorf("error banning user: %w", err)
	}
	return nil
}

// Unban unbans the user from the provided room.
// To not specify a reason, pass an empty string to the reason parameter.
func (c *Client) Unban(roomID matrix.RoomID, userID matrix.UserID, reason string) error {
	param := struct {
		UserID matrix.UserID `json:"user_id"`
		Reason string        `json:"reason,omitempty"`
	}{userID, reason}

	err := c.Request(
		"POST", EndpointRoomUnban(roomID), nil,
		httputil.WithToken(), httputil.WithJSONBody(param),
	)
	if err != nil {
		return fmt.Errorf("error unbanning user: %w", err)
	}
	return nil
}
