package event

import "github.com/chanbakjsd/gotrix/matrix"

// PushRulesEvent is an event that describes all push rules for this user.
type PushRulesEvent struct {
	EventInfo `json:"-"`
	// Global is the global ruleset.
	Global matrix.PushRuleset `json:"global"`
}

// PushNotifyMessage returns true if the message should be notified by the ruleset. Currently, only
// the message body is matched.
func PushNotifyMessage(ruleset matrix.PushRuleset, message *RoomMessageEvent) (matrix.PushRule, bool) {
	rule, ok := ruleset.Override.EventMatch(map[string]string{
		"content.body":           message.Body,
		"content.formatted_body": message.FormattedBody,
	})
	if ok {
		return rule, true
	}

	rule, ok = ruleset.Content.PatternMatch(message.Body)
	if ok {
		return rule, true
	}

	return matrix.PushRule{}, false
}
