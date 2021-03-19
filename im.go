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
		return nameEvent.Name, nil
	}

	// Step 2: Check for m.room.canonical_alias state event.
	e, _ = c.RoomState(roomID, event.TypeRoomCanonicalAlias, "")
	if e != nil {
		aliasEvent := e.(event.RoomCanonicalAliasEvent)
		if aliasEvent.Alias != "" {
			return aliasEvent.Alias, nil
		}
	}

	// TODO: This should use the m.heroes field instead of an arbitrary list of members like we are doing now.
	events, err := c.RoomStates(roomID, event.TypeRoomMember)
	if err != nil {
		return "", err
	}

	heroes := make([]matrix.UserID, 0, 5)
	for _, v := range events {
		if len(heroes) == 5 {
			break
		}

		e := v.(event.RoomMemberEvent)
		heroes = append(heroes, e.UserID)
	}

	names := make([]string, 0, len(heroes))
	for _, v := range heroes {
		name, err := c.MemberName(roomID, v)
		if err != nil {
			return "", err
		}
		names = append(names, name)
	}

	if len(events) > len(heroes) {
		names = append(names, " and "+strconv.Itoa(len(events)-len(heroes))+" others")
	}
	return strings.Join(names, ", "), nil
}
