package event

import (
	"encoding/json"
	"errors"
)

// ErrUnknownEventType represents an error where the event type is unknown and therefore
// cannot be mapped to its concrete type.
var ErrUnknownEventType = errors.New("unknown event type")

// Parse returns the concrete type of the provided event's type and sets its Event
// field to the provided event.
func (e Event) Parse() (Content, error) {
	switch e.Type {
	case TypeRoomCanonicalAlias:
		c := RoomCanonicalAliasEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomCreate:
		c := RoomCreateEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomJoinRules:
		c := RoomJoinRulesEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomMember:
		c := RoomMemberEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomPowerLevels:
		c := RoomPowerLevelsEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomRedaction:
		c := RoomRedactionEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomMessage:
		c := RoomMessageEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	}

	return nil, ErrUnknownEventType
}
