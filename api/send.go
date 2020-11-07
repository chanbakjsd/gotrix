package api

import (
	"errors"
	"net/url"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// ErrNoSendPerm represents an error where the user does not have permission to send an event to a room.
var ErrNoSendPerm = errors.New("no permission to send event")

// RoomStateSendArg represents all required arguments to (*Client).RoomStateSend.
type RoomStateSendArg struct {
	Type     event.Type
	StateKey string
	Content  interface{}
}

// RoomStateSend sends the provided state event to the provided room ID.
//
// It implements the `PUT _matrix/client/r0/rooms/{roomId}/state/{eventType}/{stateKey}` endpoint.
func (c *Client) RoomStateSend(roomID matrix.RoomID, event RoomStateSendArg) (matrix.EventID, error) {
	path := "_matrix/client/r0/rooms/" + url.PathEscape(string(roomID)) + "/state/" +
		url.PathEscape(string(event.Type)) + "/" + url.PathEscape(event.StateKey)
	var resp struct {
		EventID matrix.EventID `json:"event_id"`
	}

	err := c.Request("PUT", path, &resp, httputil.WithToken(), httputil.WithBody(event.Content))
	return resp.EventID, matrix.MapAPIError(err, matrix.ErrorMap{
		matrix.CodeForbidden: ErrNoSendPerm,
	})
}

// RoomEventSend sends the provided one-off event to the provided room ID.
//
// It implements the `PUT _matrix/client/r0/rooms/{roomId}/send/{eventType}/{txnId}` endpoint.
func (c *Client) RoomEventSend(roomID matrix.RoomID, eventType event.Type, body interface{}) (matrix.EventID, error) {
	path := "_matrix/client/r0/rooms/" + url.PathEscape(string(roomID)) + "/send/" +
		url.PathEscape(string(eventType)) + "/" + NextTransactionID()
	var resp struct {
		EventID matrix.EventID `json:"event_id"`
	}

	err := c.Request("PUT", path, &resp, httputil.WithToken(), httputil.WithBody(body))
	return resp.EventID, matrix.MapAPIError(err, matrix.ErrorMap{
		matrix.CodeForbidden: ErrNoSendPerm,
	})
}

// RoomEventRedact redacts a room event as specified by the room ID and event ID.
// A user can redact events they sent out or other people's event provided they have the power level to.
//
// It implements the `PUT _matrix/client/r0/rooms/{roomId}/redact/{eventId}/{txnId}` endpoint.
func (c *Client) RoomEventRedact(roomID matrix.RoomID, eventID matrix.EventID, reason string) (matrix.EventID, error) {
	path := "_matrix/client/r0/rooms/" + url.PathEscape(string(roomID)) + "/redact/" +
		url.PathEscape(string(eventID)) + "/" + NextTransactionID()
	var resp struct {
		EventID matrix.EventID `json:"event_id"`
	}
	var err error

	if reason == "" {
		err = c.Request("PUT", path, &resp, httputil.WithToken())
	} else {
		err = c.Request("PUT", path, &resp, httputil.WithToken(), httputil.WithBody(struct {
			Reason string `json:"reason"`
		}{
			Reason: reason,
		}))
	}
	return resp.EventID, matrix.MapAPIError(err, matrix.ErrorMap{
		matrix.CodeForbidden: ErrNoSendPerm,
	})
}
