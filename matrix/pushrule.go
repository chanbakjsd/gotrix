package matrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// PushRuleID is the rule ID for a push rule.
type PushRuleID string

const (
	MasterRuleID                = ".m.rule.master"
	SuppressNoticesRuleID       = ".m.rule.suppress_notices"
	InviteForMeRuleID           = ".m.rule.invite_for_me"
	MemberEventRuleID           = ".m.rule.member_event"
	ContainsDisplayNameRuleID   = ".m.rule.contains_display_name"
	TombstoneRuleID             = ".m.rule.tombstone"
	RoomNotificationRuleID      = ".m.rule.roomnotif"
	ContainsUsernameRuleID      = ".m.rule.contains_user_name"
	CallRuleID                  = ".m.rule.call"
	EncryptedRoomOneToOneRuleID = ".m.rule.encrypted_room_one_to_one"
	RoomOneToOneRuleID          = ".m.rule.room_one_to_one"
	MessageRuleID               = ".m.rule.message"
	EncryptedRuleID             = ".m.rule.encrypted"
)

// IsServerDefault returns true if the rule ID identifies a server-default rule.
func (id PushRuleID) IsServerDefault() bool {
	return strings.HasPrefix(string(id), ".")
}

// PushRuleset describes the global push ruleset inside a PushRulesEvent.
type PushRuleset struct {
	// Override rules are checked first. They're user-configured.
	Override PushRules `json:"override"`
	// Content rules configure behavior for unencrypted messages that match certain patterns.
	Content PushRules `json:"content"`
	// Room rules configure behavior for all messages within a given room. Their rule IDs are always
	// the ID of the rooms that they affect.
	Room PushRules `json:"room"`
	// Sender rules configure behavior for all messages from a specific Matrix user ID. Their rule
	// IDs are always the ID of those users.
	Sender PushRules `json:"sender"`
	// Underride rules are identical to Override rules but have a lower priority than all other
	// types of rules.
	Underride PushRules `json:"underride"`
}

// PushRules is a list of push rules.
type PushRules []PushRule

// Rule searches for the rule with the given ID.
func (r PushRules) Rule(id PushRuleID) (PushRule, bool) {
	for _, rule := range r {
		if rule.RuleID == id {
			return rule, true
		}
	}
	return PushRule{}, false
}

// PatternMatches returns true if any of the rules have a pattern that matches value. It behaves
// similarly to EventMatch, except the Pattern field is used instead.
func (r PushRules) PatternMatch(str string) (PushRule, bool) {
	for _, rule := range r {
		if !rule.Enabled {
			continue
		}

		if rule.Pattern.Matches(str) {
			return rule, true
		}
	}

	return PushRule{}, false
}

// EventMatch finds an enabled rule with its Kind being EventMatch and its key being any of the
// given keys in the map. It then returns the rule if its pattern matches the value of that key
// along with a true boolean. A zero-value PushRule and false are returned for any other case.
func (r PushRules) EventMatch(matchers map[string]string) (PushRule, bool) {
	for _, rule := range r {
		if !rule.Enabled {
			continue
		}

		for _, condition := range rule.Conditions {
			if condition.Kind != EventMatchCondition {
				continue
			}

			str, ok := matchers[condition.Key]
			if !ok {
				continue
			}

			if condition.Pattern.Matches(str) {
				return rule, true
			}
		}
	}

	return PushRule{}, false
}

// PushRule describes a push rule.
type PushRule struct {
	// Actions are the actions to perform when this rule is matched.
	Actions PushActions `json:"actions"`
	// Conditions are the conditions that must hold true for an event tin order for a rule to be
	// applied to an event. This only applies to Underride and Override rules.
	Conditions []PushCondition `json:"conditions,omitempty"`
	// Default (required) is true if this is a default rule or has been set explicitly.
	Default bool `json:"default"`
	// Enabled (required) is true if the rule is enabled.
	Enabled bool `json:"enabled"`
	// Pattern is the glob-style pattern to match against. This only applies to Content rules.
	Pattern PushPattern `json:"pattern"`
	// RuleID is the ID of this rule.
	RuleID PushRuleID `json:"rule_id"`
}

// PushActions describes the list of actions associated with push rules.
type PushActions struct {
	Action PushAction
	Tweaks map[PushActionTweak]json.RawMessage
}

type pushActionTweak struct {
	SetTweak PushActionTweak `json:"set_tweak"`
	Value    json.RawMessage `json:"value,omitempty"`
}

// UnmarshalJSON unmarshals an array of push actions into a.
func (a *PushActions) UnmarshalJSON(b []byte) error {
	*a = PushActions{}

	var values []json.RawMessage
	if err := json.Unmarshal(b, &values); err != nil {
		return err
	}

	if len(values) == 0 {
		return errors.New("actions have no value")
	}

	if err := json.Unmarshal(values[0], &a.Action); err != nil {
		return fmt.Errorf("cannot unmarshal first value: %w", err)
	}

	if len(values) > 1 {
		a.Tweaks = make(map[PushActionTweak]json.RawMessage, len(values[1:]))

		for _, tweakJSON := range values[1:] {
			var tweak pushActionTweak

			if err := json.Unmarshal(tweakJSON, &tweak); err != nil {
				return fmt.Errorf("cannot unmarshal tweak: %w", err)
			}

			a.Tweaks[tweak.SetTweak] = tweak.Value
		}
	}

	return nil
}

// MarshalJSON marshals PushActions into an array.
func (a PushActions) MarshalJSON() ([]byte, error) {
	values := make([]json.RawMessage, 0, 1+len(a.Tweaks))

	action, err := json.Marshal(a.Action)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal action: %w", err)
	}
	values = append(values, action)

	for key, value := range a.Tweaks {
		b, err := json.Marshal(pushActionTweak{key, value})
		if err != nil {
			return nil, fmt.Errorf("cannot marshal tweak %q: %w", key, err)
		}
		values = append(values, b)
	}

	return json.Marshal(values)
}

// SetTweak sets a tweak into the Tweaks map. If the map is nil, then it is initialized. If value is
// nil, then the value property is omitted from the SetTweak object. To delete tweaks, use the
// delete() built-in. If value is not valid JSON, then an error is returned.
func (a *PushActions) SetTweak(tweak PushActionTweak, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if a.Tweaks == nil {
		a.Tweaks = make(map[PushActionTweak]json.RawMessage, 1)
	}

	a.Tweaks[tweak] = b
	return nil
}

// UnmarshalTweak unmarshals the value with the given tweak key into ptr. False is returned if the
// JSON is invalid, missing or the tweak isn't there.
func (a PushActions) UnmarshalTweak(tweak PushActionTweak, ptr interface{}) bool {
	raw := a.Tweaks[HighlightActionTweak]
	if raw == nil {
		return false
	}

	return json.Unmarshal(raw, ptr) == nil
}

// Highlight returns true if the PushActions has a HighlightAction tweak with a true value. It
// handles certain edge cases.
func (a PushActions) Highlight() bool {
	raw, ok := a.Tweaks[HighlightActionTweak]
	if !ok {
		return false
	}

	if raw == nil {
		return true
	}

	var hl bool
	json.Unmarshal(raw, &hl)
	return hl
}

// PushAction affects if and how a notification is delivered for a matching event.
type PushAction string

const (
	// NotifyAction causes each matching event to generate a notification.
	NotifyAction = "notify"
	// DontNotifyAction prevents each matching event from generating a notification.
	DontNotifyAction = "dont_notify"
	// CoalesceAction enables notifications for matching events but activates homeserver-specific
	// behavior to intelligently coalesce multiple events into a single notification.
	CoalesceAction = "coalesce"
)

// PushActionTweak describes an entry key in the Tweaks dictionary that is sent in the notification
// request to the Push Gateway.
type PushActionTweak string

const (
	// SoundActionTweak maps to the sound to be played when this notification arrives. Default means
	// to play the default sound, but could also be something else like vibration.
	SoundActionTweak PushActionTweak = "sound"
	// HighlightActionTweak is a boolean representing whether or not this message should be
	// highlighted in the UI. If this tweak is given with no value, then it's true; if it's missing,
	// then it's false.
	HighlightActionTweak PushActionTweak = "highlight"
)

type PushConditionKind string

const (
	// EventMatchCondition is a glob pattern match on a field of the event. This condition has the
	// fields Key and Pattern.
	EventMatchCondition PushConditionKind = "event_match"
	// ContainsDisplayNameCondition matches unencrypted messages where content.body contains the
	// owner’s display name in that room. This condition has no fields.
	ContainsDisplayNameCondition PushConditionKind = "contains_display_name"
	// RoomMemberCountCondition matches the current number of members in the room. This condition
	// has the Is field. Use IsCmp to compare the count.
	RoomMemberCountCondition PushConditionKind = "room_member_count"
	// SenderNotificationPermissionCondition takes into account the current power levels in the
	// room, ensuring the sender of the event has high enough power to trigger the notification.
	// This condition has the Key field.
	SenderNotificationPermissionCondition PushConditionKind = "sender_notification_permission"
)

type PushCondition struct {
	// Is is required for room_member_count conditions. It is a decimal integer optionally prefixed
	// by one of "==", "<", ">", ">=" or "<=". If no prefix, then "==" is implied. Use the IsCmp
	// method.
	Is string `json:"is,omitempty"`
	// Key is required for:
	//
	//    - EventMatchCondition: it's a dot-separated field of the event to match.
	//    - SenderNotificationPermissionCondition: it's the field in the power level event that the
	//      user needs a minimum power level for.
	Key string `json:"key,omitempty"`
	// Kind is the kind of condition to apply.
	Kind PushConditionKind `json:"kind"`
	// Pattern is required for event_match conditions (Kind). It is the glob-style pattern to match
	// against. Patterns with no special glob characters should be treated as having asterisks
	// prepended and appended when testing the condition.
	Pattern PushPattern `json:"pattern,omitempty"`
}

// IsCmp parses the Is string and compares it with num.
func (c PushCondition) IsCmp(num int) bool {
	value := c.Is
	var valuePrefix string

	for _, prefix := range []string{"==", "<", ">", ">=", "<="} {
		if strings.HasPrefix(value, prefix) {
			value = strings.TrimPrefix(value, prefix)
			valuePrefix = prefix
			break
		}
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return false
	}

	switch valuePrefix {
	case "<":
		return i < num
	case "<=":
		return i <= num
	case ">":
		return i > num
	case ">=":
		return i >= num
	case "==", "":
		return i == num
	default:
		panic("unreachable")
	}
}

// PushPattern describes a glob string used to match push rule patterns.
type PushPattern string

// Matches returns true if str matches the Pattern string.
func (p PushPattern) Matches(str string) bool {
	pattern := string(p)
	if !strings.Contains(pattern, "*") {
		pattern = "*" + pattern + "*"
	}

	// "the push rules stuff is known to be terrible"
	matched, _ := filepath.Match(pattern, str)
	return matched
}
