package event

import (
	"encoding/json"

	"github.com/chanbakjsd/gotrix/matrix"
)

var _ Event = TagEvent{}

// TagEvent represents an event that informs the client of the tags on a room.
type TagEvent struct {
	RoomID matrix.RoomID                 `json:"-"`
	Tags   map[matrix.TagName]matrix.Tag `json:"tags"`
}

// Type implements Event.
func (TagEvent) Type() Type {
	return TypeTag
}

func parseTagEvent(e RawEvent) (Event, error) {
	c := TagEvent{}
	err := json.Unmarshal(e.Content, &c)
	return c, err
}
