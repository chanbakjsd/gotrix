package event

import (
	"encoding/json"

	"github.com/chanbakjsd/gotrix/matrix"
)

var _ Event = &TypingEvent{}

// TypingEvent is an event that updates the list of users that are typing.
type TypingEvent struct {
	EventInfo `json:"-"`

	UserID []matrix.UserID `json:"user_ids"`
	RoomID matrix.RoomID   `json:"-"`
}

func parseTypingEvent(r RawEvent, content json.RawMessage) (Event, error) {
	var tmp struct {
		EventInfo
		RoomID matrix.RoomID `json:"room_id"`
	}

	err := json.Unmarshal(r, &tmp)
	if err != nil {
		return nil, err
	}

	var v TypingEvent
	err = json.Unmarshal(content, &v)
	if err != nil {
		return nil, err
	}

	v.EventInfo = tmp.EventInfo
	v.RoomID = tmp.RoomID
	v.Raw = r

	return &v, err
}
