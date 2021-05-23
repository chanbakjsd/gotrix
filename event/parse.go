package event

import (
	"encoding/json"
	"errors"
	"reflect"
)

// ErrUnknownEventType represents an error where the event type is unknown and therefore
// cannot be mapped to its concrete type.
var ErrUnknownEventType = errors.New("unknown event type")

// Register registers a parser for the provided event type.
func Register(eventType Type, p func(RawEvent) (Event, error)) {
	parser[eventType] = p
}

// Parse returns the concrete type of the provided event's type and sets its Event
// field to the provided event.
func (e RawEvent) Parse() (Event, error) {
	p, ok := parser[e.Type]
	if !ok {
		return nil, ErrUnknownEventType
	}

	parsed, err := p(e)
	if err != nil {
		return nil, err
	}
	if parsed.Type() != e.Type {
		return nil, ErrUnknownEventType
	}
	return parsed, nil
}

type eventWithRoomEventInfo interface {
	Event
	SetRoomEventInfo(RoomEventInfo)
}

func eventParse(zeroValue func() Event) func(RawEvent) (Event, error) {
	return func(e RawEvent) (Event, error) {
		v := zeroValue()
		err := json.Unmarshal(e.Content, &v)
		w := reflect.Indirect(reflect.ValueOf(v))
		return w.Interface().(Event), err
	}
}

func roomEventParse(zeroValue func() eventWithRoomEventInfo) func(RawEvent) (Event, error) {
	return func(e RawEvent) (Event, error) {
		v := zeroValue()
		v.SetRoomEventInfo(e.toRoomEventInfo())
		err := json.Unmarshal(e.Content, &v)
		w := reflect.Indirect(reflect.ValueOf(v))
		return w.Interface().(Event), err
	}
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
