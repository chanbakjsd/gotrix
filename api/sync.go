package api

import (
	"strconv"

	"github.com/chanbakjsd/gomatrix/api/httputil"
	"github.com/chanbakjsd/gomatrix/event"
	"github.com/chanbakjsd/gomatrix/matrix"
)

// SyncArg represents all possible arguments that can be provided to a sync API call.
type SyncArg struct {
	Filter      string          `json:"filter,omitempty"`
	Since       string          `json:"since,omitempty"`
	FullState   bool            `json:"full_state,omitempty"`
	NewPresence matrix.Presence `json:"set_presence,omitempty"`
	Timeout     int             `json:"timeout,omitempty"`
}

// SyncResponse represents the response of a sync API call consisting of every info the server
// updates the client on.
type SyncResponse struct {
	NextBatch              string          `json:"next_batch"`
	Presence               SyncEvents      `json:"presence,omitempty"`
	AccountData            SyncEvents      `json:"account_data,omitempty"`
	Rooms                  SyncRoomEvents  `json:"rooms,omitempty"`
	ToDevice               SyncEvents      `json:"to_device,omitempty"`
	DeviceLists            SyncDeviceLists `json:"device_lists,omitempty"`
	DeviceOneTimeKeysCount map[string]int  `json:"device_one_time_keys_count,omitempty"`
}

// SyncRoomEvents consists of events that are tied to specific rooms (like messages and typing
// notifications).
type SyncRoomEvents struct {
	Joined  SyncJoinedRoomEvents  `json:"join,omitempty"`
	Invited SyncInvitedRoomEvents `json:"invite,omitempty"`
	Left    SyncLeftRoomEvents    `json:"leave,omitempty"`
}

// SyncJoinedRoomEvents consists of events that are tied to joined rooms (rooms the user is in).
type SyncJoinedRoomEvents struct {
	Summary     SyncRoomSummary `json:"summary,omitempty"`
	State       SyncEvents      `json:"state,omitempty"`
	Timeline    SyncTimeline    `json:"timeline,omitempty"`
	Ephemeral   SyncEvents      `json:"ephemeral,omitempty"`
	AccountData SyncEvents      `json:"account_data,omitempty"`
	UnreadCount struct {
		Highlight    int `json:"highlight_count,omitempty"`
		Notification int `json:"notification_count,omitempty"`
	} `json:"unread_notifications,omitempty"`
}

// SyncInvitedRoomEvents consists of events that are tied to rooms that the client is invited to.
type SyncInvitedRoomEvents struct {
	State struct {
		Events []event.StrippedEvent `json:"events,omitempty"`
	} `json:"state,omitempty"`
}

// SyncLeftRoomEvents consists of events that are tied to rooms that the user has left.
type SyncLeftRoomEvents struct {
	State       SyncEvents   `json:"state,omitempty"`
	Timeline    SyncTimeline `json:"timeline,omitempty"`
	AccountData SyncEvents   `json:"account_data,omitempty"`
}

// SyncRoomSummary consists of data that the client may need to render a room correctly.
//
// Heroes are users that are allowed to set a name/canonical alias to a room.
type SyncRoomSummary struct {
	Heroes       []string `json:"m.heroes,omitempty"`
	JoinedCount  int      `json:"m.joined_member_count,omitempty"`
	InvitedCount int      `json:"m.invited_member_count,omitempty"`
}

// SyncTimeline consists of a timeline of events.
//
// If limited is true, the query is limited by the filter limits set and the client should
// query if desired.
//
// PreviousBatch can be used as an ID to index into previous timeline events.
type SyncTimeline struct {
	Events        []event.Event `json:"events,omitempty"`
	Limited       bool          `json:"limited,omitempty"`
	PreviousBatch string        `json:"prev_batch"`
}

// SyncEvents are a list of events.
type SyncEvents struct {
	Events []event.Event `json:"events,omitempty"`
}

// SyncDeviceLists is a list of users who has their encryption keys changed (added or modified)
// or deleted (Left).
type SyncDeviceLists struct {
	Changed []matrix.UserID
	Left    []matrix.UserID
}

// Sync requests the latest state changes from the server.
//
// It implements the `GET _matrix/client/r0/sync` endpoint.
func (c *Client) Sync(req SyncArg) (*SyncResponse, error) {
	var resp *SyncResponse
	args := make(map[string]string)
	if req.Filter != "" {
		args["filter"] = req.Filter
	}
	if req.Since != "" {
		args["since"] = req.Since
	}
	if req.FullState {
		args["full_state"] = "true"
	}
	if req.NewPresence != "" {
		args["set_presence"] = string(req.NewPresence)
	}
	if req.Timeout != 0 {
		args["timeout"] = strconv.Itoa(req.Timeout)
	}
	err := c.Request(
		"GET", "_matrix/client/r0/sync", resp,
		httputil.WithToken(), httputil.WithQuery(args),
	)
	return resp, err
}