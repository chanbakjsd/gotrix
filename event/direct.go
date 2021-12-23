package event

import (
	"encoding/json"

	"github.com/chanbakjsd/gotrix/matrix"
)

var _ Event = &DirectEvent{}

// DirectEvent is an event that lists all the DM channels the user is in.
// It is saved in AccountData.
type DirectEvent struct {
	EventInfo `json:"-"`
	Rooms     map[matrix.UserID][]matrix.RoomID
}

func parseDirectEvent(r RawEvent, content json.RawMessage) (Event, error) {
	var v DirectEvent
	err := json.Unmarshal(r, &v.EventInfo)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &v.Rooms)
	if err != nil {
		return nil, err
	}

	v.Raw = r
	return &v, nil
}

// MarshalJSON marshals the internal list of rooms inside DirectEvent to be consistent with other
// events as they only include content in the marshalled JSON.
func (d DirectEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Rooms)
}
