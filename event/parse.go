package event

import (
	"encoding/json"
	"fmt"
)

// UnknownEventTypeError represents an error where the event type is unknown and therefore
// cannot be mapped to its concrete type.
type UnknownEventTypeError struct {
	Found Type
}

func (e UnknownEventTypeError) Error() string {
	return fmt.Sprintf("unknown event type: %s", e.Found)
}

// Register registers a parser for the provided event type.
// The parser is passed the full raw event and its content field.
func Register(eventType Type, p func(RawEvent, json.RawMessage) (Event, error)) {
	parser[eventType] = p
}

// RegisterDefault registers a parser for the provided event type.
// It automatically fills in the applicable EventInfo and passes only the content to the parser
// function.
func RegisterDefault(eventType Type, p func(json.RawMessage) (Event, error)) {
	parser[eventType] = func(r RawEvent, content json.RawMessage) (Event, error) {
		v, err := p(content)
		if err != nil {
			return nil, err
		}

		return fillInfo(r, v)
	}
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

	p, ok := parser[ev.Type]
	if !ok {
		return nil, UnknownEventTypeError{
			Found: ev.Type,
		}
	}

	// Copy RawEvent before passing it to the individual parsers as most store
	// it in the Raw field literally.
	r = append([]byte(nil), r...)
	concrete, err := p(r, ev.Content)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling event of type %s: %w", ev.Type, err)
	}

	if concrete.Info().Type != ev.Type {
		return nil, fmt.Errorf("error unmarshalling event of type %s: got type %s from registered parser", ev.Type, concrete.Info().Type)
	}

	return concrete, nil
}

// fillInfo is a helper function that backfills EventInfo/RoomEventInfo/StateEventInfo into the provided event.
func fillInfo(raw RawEvent, v Event) (Event, error) {
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

// defaultParse parses the content into the event returned by the zeroValue func, assuming that it
// creates a pointer. It then fills in the info.
func defaultParse(zeroValue func() Event) func(RawEvent, json.RawMessage) (Event, error) {
	return func(raw RawEvent, content json.RawMessage) (Event, error) {
		v := zeroValue()
		err := json.Unmarshal(content, v)
		if err != nil {
			return nil, err
		}

		return fillInfo(raw, v)
	}
}
