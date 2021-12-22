package event

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrUnknownEventType represents an error where the event type is unknown and therefore
// cannot be mapped to its concrete type.
var ErrUnknownEventType = errors.New("unknown event type")

// Register registers a parser for the provided event type.
// The parser is passed the full raw event and its content field.
func Register(eventType Type, p func(RawEvent, json.RawMessage) (Event, error)) {
	parser[eventType] = p
}

// Parse returns the concrete type of the provided event's type.
func Parse(r RawEvent) (Event, error) {
	type event struct {
		Type    Type            `json:"type"`
		Content json.RawMessage `json:"content"`
	}

	var ev event
	err := json.Unmarshal(r, &ev)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling event: %w", err)
	}

	if _, ok := parser[ev.Type]; !ok {
		return nil, ErrUnknownEventType
	}

	concrete, err := parser[ev.Type](r, ev.Content)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling event of type %s: %w", ev.Type, err)
	}

	if concrete.Info().Type != ev.Type {
		return nil, fmt.Errorf("error unmarshalling event of type %s: got type %s from registered parser", ev.Type, concrete.Info().Type)
	}

	return concrete, nil
}

func defaultParse(zeroValue func() Event) func(RawEvent, json.RawMessage) (Event, error) {
	return func(raw RawEvent, content json.RawMessage) (Event, error) {
		v := zeroValue()
		err := json.Unmarshal(content, v)
		if err != nil {
			return nil, err
		}

		switch e := v.(type) {
		case StateEvent:
			err := json.Unmarshal(raw, e.StateInfo())
			if err != nil {
				return nil, err
			}
		case RoomEvent:
			err := json.Unmarshal(raw, e.RoomInfo())
			if err != nil {
				return nil, err
			}
		default:
			err := json.Unmarshal(raw, e.Info())
			if err != nil {
				return nil, err
			}
		}

		v.Info().Raw = raw

		return v, nil
	}
}
