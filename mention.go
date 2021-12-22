package gotrix

import (
	"net/url"

	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// MentionUser creates the intended mention format for a user in normal body and formatted body.
func (c *Client) MentionUser(userID matrix.UserID, roomID matrix.RoomID) (body string, formatted string) {
	e, _ := c.RoomState(roomID, event.TypeRoomMember, string(userID))
	if e == nil {
		return formatMention(string(userID), string(userID))
	}

	memberEvent := e.(*event.RoomMemberEvent)
	if memberEvent.DisplayName == nil || *memberEvent.DisplayName == "" {
		return formatMention(string(userID), string(userID))
	}

	return formatMention(string(userID), *memberEvent.DisplayName)
}

// MentionRoom creates the intended mention format for a room in normal body and formatted body.
func (c *Client) MentionRoom(roomID matrix.RoomID) (body string, formatted string) {
	e, _ := c.RoomState(roomID, event.TypeRoomCanonicalAlias, "")
	if e != nil {
		aliasEvent := e.(*event.RoomCanonicalAliasEvent)
		if aliasEvent.Alias != "" {
			return formatMention(string(roomID), aliasEvent.Alias)
		}
		if len(aliasEvent.AltAlias) > 0 {
			return formatMention(string(roomID), aliasEvent.AltAlias[0])
		}
	}

	return formatMention(string(roomID), string(roomID))
}

// TODO Implement group mention

func formatMention(id, name string) (string, string) {
	formatted := "<a href='https://matrix.to/#/" + url.PathEscape(id) + "'>" + name + "</a>"
	return name, formatted
}
