package event

import (
	"encoding/json"

	"github.com/chanbakjsd/gotrix/matrix"
)

var _ Event = DirectEvent{}

// DirectEvent is an event that lists all the DM channels the user is in.
type DirectEvent map[matrix.UserID][]matrix.RoomID

// Type implements Event.
func (DirectEvent) Type() Type {
	return TypeDirect
}

// Raw returns DirectEvent in the form of a RawEvent.
func (d DirectEvent) Raw() (RawEvent, error) {
	bytes, err := json.Marshal(d)
	if err != nil {
		return RawEvent{}, err
	}

	return RawEvent{
		Type:    TypeDirect,
		Content: json.RawMessage(bytes),
	}, nil
}
