package gotrix

import (
	"errors"
	"strconv"
	"strings"

	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// ErrRoomAvatarNotFound is returned by RoomAvatar when the room avatar cannot be discovered.
var ErrRoomAvatarNotFound = errors.New("user is alone in the room and a room avatar has not been set")

// RoomAvatar retrieves the avatar of a room.
// It falls back to the profile picture of the first user to join the room that is not the current user otherwise.
// If the current user is alone, it returns nil and an error.
func (c *Client) RoomAvatar(roomID matrix.RoomID) (*matrix.URL, error) {
	e, _ := c.RoomState(roomID, event.TypeRoomAvatar, "")
	if e != nil {
		avatarEvent := e.(*event.RoomAvatarEvent)
		return &avatarEvent.URL, nil
	}

	summary, err := c.RoomSummary(roomID)
	if err != nil {
		return nil, err
	}
	if len(summary.Heroes) == 0 {
		return nil, ErrRoomAvatarNotFound
	}

	return c.MemberAvatar(roomID, summary.Heroes[0])
}

// RoomName calculates the display name of a room.
func (c *Client) RoomName(roomID matrix.RoomID) (string, error) {
	// Step 1: Check for m.room.name state event.
	e, _ := c.RoomState(roomID, event.TypeRoomName, "")
	if e != nil {
		nameEvent := e.(*event.RoomNameEvent)
		if nameEvent.Name != "" {
			return nameEvent.Name, nil
		}
	}

	// Step 2: Check for m.room.canonical_alias state event.
	e, _ = c.RoomState(roomID, event.TypeRoomCanonicalAlias, "")
	if e != nil {
		aliasEvent := e.(*event.RoomCanonicalAliasEvent)
		if aliasEvent.Alias != "" {
			return aliasEvent.Alias, nil
		}
	}

	summary, err := c.RoomSummary(roomID)
	if err != nil {
		return "", err
	}

	heroes := make([]string, 0, len(summary.Heroes))
	for k, v := range summary.Heroes {
		if k > 4 {
			break // Sane limit of 5 names displayed.
		}
		name, err := c.MemberName(roomID, v)
		if err != nil {
			return "", err
		}
		heroes = append(heroes, name)
	}

	joinAndInvited := summary.JoinedCount + summary.InvitedCount
	if len(heroes) == 0 {
		if joinAndInvited <= 1 {
			// User is alone in the room or the room is empty.
			return "Empty Room", nil
		}

		// This should never happen but if there are no heroes, the room ID is the sanest option we have.
		return string(roomID), nil
	}

	switch {
	case len(summary.Heroes) == 1:
		// Do nothing. There's no "and" to add.
	case len(summary.Heroes) >= joinAndInvited-1 && len(heroes) > 1:
		// There are only heroes in the room so just make it "and <Last Hero>".
		heroes[len(heroes)-1] = "and " + heroes[len(heroes)-1]
	default:
		// There are more than just heroes in the room.
		heroes = append(heroes, "and "+strconv.Itoa(joinAndInvited-len(summary.Heroes))+" others")
	}

	roomSummary := strings.Join(heroes, ", ")
	if len(heroes) == 2 {
		// Do "Alice and Bob" instead of "Alice, and Bob".
		roomSummary = heroes[0] + " " + heroes[1]
	}

	if joinAndInvited <= 1 {
		return "Empty Room (was " + roomSummary + ")", nil
	}
	return roomSummary, nil
}
