package event

import (
	"encoding/json"
	"errors"

	"github.com/chanbakjsd/gotrix/matrix"
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
			Event:  e,
			UserID: matrix.UserID(e.StateKey),
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
	case TypeRoomName:
		c := RoomNameEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomTopic:
		c := RoomTopicEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomAvatar:
		c := RoomAvatarEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomPinned:
		c := RoomPinnedEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeCallInvite:
		c := CallInviteEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeCallCandidates:
		c := CallCandidatesEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeCallAnswer:
		c := CallAnswerEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeCallHangup:
		c := CallHangupEvent{
			Event: e,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	}

	return nil, ErrUnknownEventType
}
