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
func (e RawEvent) Parse() (Event, error) {
	switch e.Type {
	case TypeRoomCanonicalAlias:
		c := RoomCanonicalAliasEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomCreate:
		c := RoomCreateEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomJoinRules:
		c := RoomJoinRulesEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomMember:
		c := RoomMemberEvent{
			RoomEventInfo: e.toRoomEventInfo(),
			UserID:        matrix.UserID(e.StateKey),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomPowerLevels:
		c := RoomPowerLevelsEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomRedaction:
		c := RoomRedactionEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomMessage:
		c := RoomMessageEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomName:
		c := RoomNameEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomTopic:
		c := RoomTopicEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomAvatar:
		c := RoomAvatarEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeRoomPinned:
		c := RoomPinnedEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeCallInvite:
		c := CallInviteEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeCallCandidates:
		c := CallCandidatesEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeCallAnswer:
		c := CallAnswerEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeCallHangup:
		c := CallHangupEvent{
			RoomEventInfo: e.toRoomEventInfo(),
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeTyping:
		c := TypingEvent{
			RoomID: e.RoomID,
		}
		err := json.Unmarshal(e.Content, &c)
		return c, err
	case TypeReceipt:
		c := ReceiptEvent{
			RoomID: e.RoomID,
		}
		err := json.Unmarshal(e.Content, &c.Events)
		return c, err
	}

	return nil, ErrUnknownEventType
}

// toRoomEventInfo creates a RoomEventInfo from the provided event.
func (e RawEvent) toRoomEventInfo() RoomEventInfo {
	return RoomEventInfo{
		EventID:    e.ID,
		RoomID:     e.RoomID,
		SenderID:   e.Sender,
		OriginTime: e.OriginServerTime,
	}
}
