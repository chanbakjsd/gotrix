package gotrix

import (
	"errors"

	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// MemberAvatarNotFound is returned by MemberAvatar when a member's avatar cannot be determined.
var MemberAvatarNotFound = errors.New("member avatar cannot be found")

// MemberAvatar determines the member's avatar by checking for a room member event with avatar and
// falling back to their global avatar.
// An error is returned if the avatar cannot be determined.
func (c *Client) MemberAvatar(roomID matrix.RoomID, userID matrix.UserID) (*matrix.URL, error) {
	e, err := c.RoomState(roomID, event.TypeRoomMember, string(userID))
	if err != nil {
		return nil, err
	}

	memberEvent := e.(event.RoomMemberEvent)
	if memberEvent.AvatarURL != "" {
		return &memberEvent.AvatarURL, nil
	}

	return c.Client.AvatarURL(userID)
}

// MemberName calculates the display name of a member.
// Note that a user joining might invalidate some names if they share the same display name as disambiguation
// will become necessary.
//
// Use the Client.MemberNames variant when generating member name for multiple users to reduce duplicate work.
func (c *Client) MemberName(roomID matrix.RoomID, userID matrix.UserID) (string, error) {
	names, err := c.MemberNames(roomID, []matrix.UserID{userID})
	if err != nil {
		return "", err
	}
	return names[0], nil
}

// MemberNames calculates the display name of all the users provided.
func (c *Client) MemberNames(roomID matrix.RoomID, userIDs []matrix.UserID) ([]string, error) {
	// Build the hashmap of display names to locate duplicate display names.
	dupe := make(map[string]int)
	err := c.EachRoomState(roomID, event.TypeRoomMember, func(key string, v event.StateEvent) error {
		memberEvent := v.(event.RoomMemberEvent)
		if memberEvent.DisplayName == nil {
			return nil
		}
		dupe[*memberEvent.DisplayName]++
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Start generating display names.
	result := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		// Step 1: Inspect m.room.member state event.
		e, _ := c.RoomState(roomID, event.TypeRoomMember, string(userID))
		if e == nil {
			result = append(result, string(userID))
			continue
		}

		// Step 2: If there are no display name field, use raw user ID as display name.
		memberEvent := e.(event.RoomMemberEvent)
		if memberEvent.DisplayName == nil || *memberEvent.DisplayName == "" {
			result = append(result, string(userID))
			continue
		}

		displayName := *memberEvent.DisplayName
		// Step 3: Use display name if it is unique.
		if dupe[displayName] == 1 {
			result = append(result, displayName)
			continue
		}

		// Step 4: Disambiguate if the display name is not unique.
		result = append(result, displayName+" ("+string(userID)+")")
	}

	return result, nil
}
