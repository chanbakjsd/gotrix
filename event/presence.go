package event

import (
	"encoding/json"
	"time"

	"github.com/chanbakjsd/gotrix/matrix"
)

var _ Event = &PresenceEvent{}

// PresenceEvent is an event where the presence of a user is updated.
type PresenceEvent struct {
	EventInfo `json:"-"`

	User        matrix.UserID `json:"-"`
	AvatarURL   *matrix.URL   `json:"avatar_url,omitempty"`
	DisplayName *string       `json:"displayname,omitempty"`

	// Last time since user performed some action, in ms.
	LastActiveAgo   *int            `json:"last_active_ago,omitempty"`
	Presence        matrix.Presence `json:"presence"`
	CurrentlyActive *bool           `json:"currently_active,omitempty"`
	Status          *string         `json:"status_msg,omitempty"`

	receiveTime time.Time
}

// LastActive calculates the last active time based on the time the event is parsed and the last active ago field.
// It is slightly off as the time the event is received is subject to network latency.
// It returns nil if the last active ago field is absent.
func (p PresenceEvent) LastActive() *time.Time {
	if p.LastActiveAgo == nil {
		return nil
	}
	lastActive := p.receiveTime.Add(-time.Duration(*p.LastActiveAgo) * time.Millisecond)
	return &lastActive
}

func parsePresenceEvent(r RawEvent, content json.RawMessage) (Event, error) {
	receiveTime := time.Now()
	var fullEvent struct {
		EventInfo
		Content PresenceEvent `json:"content"`
		Sender  matrix.UserID `json:"sender"`
	}

	err := json.Unmarshal(r, &fullEvent)
	if err != nil {
		return nil, err
	}

	v := fullEvent.Content
	v.EventInfo = fullEvent.EventInfo
	v.User = fullEvent.Sender
	v.receiveTime = receiveTime
	v.Raw = r

	return &v, nil
}
