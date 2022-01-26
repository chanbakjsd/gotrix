package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// RoomCreateArg represents all arguments to (*Client).RoomCreate.
type RoomCreateArg struct {
	Visibility       RoomVisibility     `json:"visibility,omitempty"`
	AliasName        string             `json:"room_alias_name,omitempty"` // Desired local part of room alias
	Name             string             `json:"name,omitempty"`            // The displayed name
	Topic            string             `json:"topic,omitempty"`           // The displayed topic
	Invite           []matrix.UserID    `json:"invite,omitempty"`          // A list of users to invite
	ThirdpartyInvite []RoomCreateInvite `json:"invite_3pid,omitempty"`     // List of third party invites
	Version          string             `json:"room_version,omitempty"`    // Room Version
	InitialState     []event.Event      `json:"initial_state,omitempty"`   // Initial State
	Preset           RoomPreset         `json:"preset,omitempty"`          // The preset to use for permissions.
	IsDirectMessage  bool               `json:"is_direct,omitempty"`       // True if this should be a Direct Message
	// Extra keys to add to the RoomCreate event.
	CreationContent    map[string]interface{}      `json:"creation_content,omitempty"`
	PowerLevelOverride *event.RoomPowerLevelsEvent `json:"power_level_content_override,omitempty"`
}

// RoomVisibility represents the initial visibility of the room.
type RoomVisibility string

const (
	// RoomPublic will publish the room into the published room list.
	RoomPublic RoomVisibility = "public"
	// RoomPrivate will NOT publish the room into the published room list.
	RoomPrivate RoomVisibility = "private"
)

// RoomPreset is the preset of room settings that can be used for sane defaults.
type RoomPreset string

const (
	// PresetTrustedPrivateChat gives everyone admin practically.
	PresetTrustedPrivateChat RoomPreset = "trusted_private_chat"
	// PresetPrivateChat makes the room invite-only.
	PresetPrivateChat RoomPreset = "private_chat"
	// PresetPublicChat mandates account creation to view the room.
	PresetPublicChat RoomPreset = "public_chat"
)

// RoomCreateInvite represents information to identify a third party user to invite.
type RoomCreateInvite struct {
	IdentityServer string `json:"id_server"`       // Identity server to lookup on
	AccessToken    string `json:"id_access_token"` // Access token previously registered with identity server
	Medium         string `json:"medium"`          // The kind being looked up like "email".
	Address        string `json:"address"`         // The third party identifier to look up.
}

// RoomCreate creates the room with the provided arguments.
//
// Note: For backwards compatibility reasons, the default setting for history visibility is such that users
// can view all messages as long as they were a member at some point. This might be undesirable and the client
// should prompt the user to decide, sending RoomHistoryVisibilityEvent as necessary.
func (c *Client) RoomCreate(arg RoomCreateArg) (matrix.RoomID, error) {
	resp := &struct {
		RoomID matrix.RoomID `json:"room_id"`
	}{}
	err := c.Request(
		"POST", c.Endpoints.RoomCreate(), resp,
		httputil.WithToken(), httputil.WithJSONBody(arg),
	)
	if err != nil {
		return "", fmt.Errorf("error creating room: %w", err)
	}
	return resp.RoomID, nil
}
