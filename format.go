package gotrix

import (
	"fmt"
	"html"
	"io"
	"net/url"
	"strings"

	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// FormatSpoiler creates the intended spoiler format for a spoiler message.
// Adhering to the spec, the spoiler text is uploaded to MXC as a plaintext to be included in the
// body. This function should therefore not be used in encrypted rooms to prevent leaks.
func (c *Client) FormatSpoiler(reason string, spoilerText string) (body string, formatted string, _ error) {
	url, err := c.MediaUpload("text/plain", "spoiler.txt", io.NopCloser(strings.NewReader(spoilerText)))
	if err != nil {
		return "", "", err
	}
	if reason == "" {
		return "[Spoiler](" + string(url) + ")",
			"<span data-mx-spoiler>" + spoilerText + "</span>",
			nil
	}

	return fmt.Sprintf("[Spoiler for %s](%s)", reason, url),
		fmt.Sprintf("<span data-mx-spoiler='%s'>%s</span>", html.EscapeString(reason), spoilerText),
		nil
}

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
