package gotrix

import (
	"strconv"
	"strings"

	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// MemberName calculates the display name of a member.
func (c *Client) MemberName(roomID matrix.RoomID, userID matrix.UserID) (string, error) {
	// Step 1: Inspect m.room.member state event.
	e, _ := c.RoomState(roomID, event.TypeRoomMember, string(userID))
	if e == nil {
		return string(userID), nil
	}

	memberEvent := e.(event.RoomMemberEvent)
	if memberEvent.DisplayName == nil || *memberEvent.DisplayName == "" {
		return string(userID), nil
	}

	// TODO: Check for the need to disambiguate as requested by the spec 13.2.2.3.
	return *memberEvent.DisplayName, nil
}

// RoomName calculates the display name of a room.
func (c *Client) RoomName(roomID matrix.RoomID) (string, error) {
	// Step 1: Check for m.room.name state event.
	e, _ := c.RoomState(roomID, event.TypeRoomName, "")
	if e != nil {
		nameEvent := e.(event.RoomNameEvent)
		if nameEvent.Name != "" {
			return nameEvent.Name, nil
		}
	}

	// Step 2: Check for m.room.canonical_alias state event.
	e, _ = c.RoomState(roomID, event.TypeRoomCanonicalAlias, "")
	if e != nil {
		aliasEvent := e.(event.RoomCanonicalAliasEvent)
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
