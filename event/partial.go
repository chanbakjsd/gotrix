package event

import "encoding/json"

// Partial is a partially parsed event object. Its pointer implements StateEvent but is not
// guaranteed to be a state event or a room event.
// It aims to allow API users to inspect into events that fail to unmarshal because it is of an
// unknown event type.
type Partial struct {
	StateEventInfo
	Content json.RawMessage `json:"content"`
}

// ParsePartial parses the raw event partially, leaving Content untouched and providing some fields
// common to most events exposed for inspection.
func ParsePartial(raw RawEvent) (*Partial, error) {
	var p Partial
	err := json.Unmarshal(raw, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Raw returns the raw form of Partial by marshalling it.
func (p Partial) Raw() (RawEvent, error) {
	return json.Marshal(p)
}
