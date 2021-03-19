package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// ErrEventNotFound is returned when the event is not found or cannot be accessed by the user.
// It is returned by (*Client).RoomEvent, (*Client).RoomState and (*Client).RoomStates.
var ErrEventNotFound = errors.New("event not found")

// ErrRoomNotFound is returned when the room is not found or cannot be accessed by the user.
// It is returned by (*Client).RoomState, (*Client).RoomStates, (*Client).RoomMembers,
// (*Client).RoomMessages, (*Client).RoomAliases and (*Client).Join.
var ErrRoomNotFound = errors.New("room not found")

// RoomEvent fetches an event from the server with the provided room ID or event ID.
//
// It implements the `GET _matrix/client/r0/rooms/{roomId}/event/{eventId}` endpoint.
func (c *Client) RoomEvent(roomID matrix.RoomID, eventID matrix.EventID) (*event.RawEvent, error) {
	path := "_matrix/client/r0/rooms/" + url.PathEscape(string(roomID)) +
		"/event/" + url.PathEscape(string(eventID))

	resp := &event.RawEvent{}
	err := c.Request("GET", path, resp, httputil.WithToken())
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching room event: %w",
			matrix.MapAPIError(err, matrix.ErrorMap{
				matrix.CodeNotFound: ErrEventNotFound,
			}),
		)
	}
	return resp, nil
}

// RoomState fetches the latest state event for the provided state in the provided room.
//
// It implements the `GET _matrix/client/r0/rooms/{roomId}/state/{eventType}/{stateKey}` endpoint.
func (c *Client) RoomState(roomID matrix.RoomID, eventType event.Type, key string) (*event.RawEvent, error) {
	path := "_matrix/client/r0/rooms/" + url.PathEscape(string(roomID)) + "/state/" +
		url.PathEscape(string(eventType)) + "/" + url.PathEscape(key)
	var content json.RawMessage
	err := c.Request("GET", path, &content, httputil.WithToken())
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching room state event: %w",
			matrix.MapAPIError(err, matrix.ErrorMap{
				matrix.CodeNotFound:  ErrEventNotFound,
				matrix.CodeForbidden: ErrRoomNotFound,
			}),
		)
	}
	return &event.RawEvent{
		Type:     eventType,
		Content:  content,
		RoomID:   roomID,
		StateKey: key,
	}, nil
}

// RoomStates fetches all the current state events of the provided room.
// If the user has left the room, it returns the state before the user leaves.
//
// It implements the `GET _matrix/client/r0/rooms/{roomId}/state` endpoint.
func (c *Client) RoomStates(roomID matrix.RoomID) (*[]event.RawEvent, error) {
	path := "_matrix/client/r0/rooms/" + url.PathEscape(string(roomID))

	resp := &[]event.RawEvent{}
	err := c.Request("GET", path, resp, httputil.WithToken())
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching room states: %w",
			matrix.MapAPIError(err, matrix.ErrorMap{
				matrix.CodeNotFound:  ErrEventNotFound,
				matrix.CodeForbidden: ErrRoomNotFound,
			}),
		)
	}
	return resp, nil
}

// RoomMemberFilter represents a filter that can be set to filter a RoomMembers request.
type RoomMemberFilter struct {
	// The pagination token to query at.
	At            string           `json:"at,omitempty"`
	Membership    event.MemberType `json:"membership,omitempty"`
	NotMembership event.MemberType `json:"not_membership,omitempty"`
}

// RoomMembers fetches the member list for a room from the homeserver.
// The returned member list is in the form of an array of RoomMember events.
//
// It implements the `GET _matrix/client/r0/rooms/{roomId}/members` endpoint.
func (c *Client) RoomMembers(roomID matrix.RoomID, filter RoomMemberFilter) (*[]event.RoomMemberEvent, error) {
	path := "_matrix/client/r0/rooms/" + url.PathEscape(string(roomID))
	var resp struct {
		Chunk *[]event.RoomMemberEvent `json:"chunk,omitempty"`
	}

	arg := make(map[string]string)
	if filter.At != "" {
		arg["at"] = filter.At
	}
	if filter.Membership != "" {
		arg["membership"] = string(filter.Membership)
	}
	if filter.NotMembership != "" {
		arg["not_membership"] = string(filter.NotMembership)
	}

	err := c.Request("GET", path, &resp, httputil.WithToken(), httputil.WithQuery(arg))
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching room members: %w",
			matrix.MapAPIError(err, matrix.ErrorMap{
				matrix.CodeForbidden: ErrRoomNotFound,
			}),
		)
	}
	return resp.Chunk, nil
}

// RoomMember represents a member in a room as returned by (*Client).RoomJoinedMembers.
type RoomMember struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// RoomJoinedMembers fetches all the joined members and return them as a map of user ID to room member.
//
// It implements the `GET _matrix/client/r0/rooms/{roomId}/joinedMembers` endpoint.
func (c *Client) RoomJoinedMembers(roomID matrix.RoomID) (*map[matrix.UserID]RoomMember, error) {
	resp := &map[matrix.UserID]RoomMember{}
	err := c.Request(
		"GET", "_matrix/client/r0/rooms/"+url.PathEscape(string(roomID))+"/joined_members", resp,
		httputil.WithToken(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching joined members: %w",
			matrix.MapAPIError(err, matrix.ErrorMap{
				matrix.CodeForbidden: ErrRoomNotFound,
			}),
		)
	}
	return resp, nil
}

// RoomMessagesDirection specifies the direction to fetch in.
type RoomMessagesDirection string

const (
	// RoomMessagesForward fetches newer messages.
	RoomMessagesForward RoomMessagesDirection = "f"
	// RoomMessagesBackward fetches older messages.
	RoomMessagesBackward RoomMessagesDirection = "b"
)

// RoomMessagesQuery represents the query the client should send to the homeserver.
type RoomMessagesQuery struct {
	From      string                // Required
	Direction RoomMessagesDirection // Required
	To        string
	Limit     int
	Filter    *event.RoomEventFilter
}

// RoomMessagesResponse represents the response to (*Client).RoomMessages.
type RoomMessagesResponse struct {
	Start string           `json:"start"` // The token pagination starts from.
	End   string           `json:"end"`   // The token pagination ends on.
	Chunk []event.RawEvent `json:"chunk"` // A list of room events.
	State []event.RawEvent `json:"state"` // A list of state events relevant to the room events.
}

// RoomMessages fetches the messages specified in the query range and return them.
//
// It implements the `GET _matrix/client/r0/rooms/{roomId}/messages` endpoint.
func (c *Client) RoomMessages(roomID matrix.RoomID, query RoomMessagesQuery) (*RoomMessagesResponse, error) {
	arg := map[string]string{
		"from": query.From,
		"dir":  string(query.Direction),
	}
	if query.To != "" {
		arg["to"] = query.To
	}
	if query.Limit != 0 {
		arg["limit"] = strconv.Itoa(query.Limit)
	}
	if query.Filter != nil {
		bytes, err := json.Marshal(query.Filter)
		if err != nil {
			return nil, fmt.Errorf("error marshalling filter: %w", err)
		}
		arg["filter"] = string(bytes)
	}

	resp := &RoomMessagesResponse{}
	err := c.Request(
		"GET", "_matrix/client/r0/rooms/"+url.PathEscape(string(roomID))+"/messages", resp,
		httputil.WithToken(), httputil.WithQuery(arg),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching room messages: %w",
			matrix.MapAPIError(err, matrix.ErrorMap{
				matrix.CodeForbidden: ErrRoomNotFound,
			}),
		)
	}
	return resp, nil
}
