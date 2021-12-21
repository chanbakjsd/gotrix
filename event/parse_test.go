package event

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/chanbakjsd/gotrix/matrix"
)

func TestSpecExamples(t *testing.T) {
	for _, v := range specExamples {
		ev, err := RawEvent(v.Code).Parse()
		if err != nil {
			t.Errorf("error parsing %s: %v", v.Name, err)
			continue
		}

		v.Expected.Info().Raw = json.RawMessage(v.Code)
		if !reflect.DeepEqual(ev, v.Expected) {
			t.Errorf("mismatch on parsing %s\nexpected: %#v\ngot: %#v", v.Name, v.Expected, ev)
			continue
		}
	}
}

var boolTrue = true

func stringPtr(a string) *string {
	return &a
}

func intPtr(a int) *int {
	return &a
}

var specExamples = []struct {
	Name     string
	Code     string
	Expected Event
}{
	{
		Name: "m.room.canonical_alias",
		Code: `
			{
				"content": {
					"alias": "#somewhere:localhost",
					"alt_aliases": [
						"#somewhere:example.org",
						"#myroom:example.com"
					]
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.canonical_alias",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: RoomCanonicalAliasEvent{
			Alias: "#somewhere:localhost",
			AltAlias: []string{
				"#somewhere:example.org",
				"#myroom:example.com",
			},
			StateEventInfo: &StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomCanonicalAlias,
					},
					ID:               "$143273582443PhrSn:example.org",
					OriginServerTime: 1432735824653,
					RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
					Sender:           "@example:example.org",
					Unsigned: UnsignedData{
						Age: 1234,
					},
				},
				StateKey: "",
			},
		},
	},
	{
		Name: "m.room.create",
		Code: `
			{
				"content": {
					"creator": "@example:example.org",
					"m.federate": true,
					"predecessor": {
						"event_id": "$something:example.org",
						"room_id": "!oldroom:example.org"
					},
					"room_version": "1"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.create",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: RoomCreateEvent{
			Creator: "@example:example.org",
			// TODO: Add predecessor field
			Federated:   &boolTrue,
			RoomVersion: stringPtr("1"),
			StateEventInfo: &StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomCreate,
					},
					ID:               "$143273582443PhrSn:example.org",
					OriginServerTime: 1432735824653,
					RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
					Sender:           "@example:example.org",
					Unsigned: UnsignedData{
						Age: 1234,
					},
				},
				StateKey: "",
			},
		},
	},
	{
		Name: "m.room.join_rules",
		Code: `
			{
				"content": {
					"join_rule": "public"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.join_rules",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: RoomJoinRulesEvent{
			JoinRule: JoinPublic,
			StateEventInfo: &StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomJoinRules,
					},
					ID:               "$143273582443PhrSn:example.org",
					OriginServerTime: 1432735824653,
					RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
					Sender:           "@example:example.org",
					Unsigned: UnsignedData{
						Age: 1234,
					},
				},
				StateKey: "",
			},
		},
	},
	{
		Name: "m.room.member 1",
		Code: `
		{
			"content": {
				"avatar_url": "mxc://example.org/SEsfnsuifSDFSSEF",
				"displayname": "Alice Margatroid",
				"membership": "join",
				"reason": "Looking for support"
			},
			"event_id": "$143273582443PhrSn:example.org",
			"origin_server_ts": 1432735824653,
			"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
			"sender": "@example:example.org",
			"state_key": "@alice:example.org",
			"type": "m.room.member",
			"unsigned": {
				"age": 1234
			}
		}
		`,
		Expected: RoomMemberEvent{
			AvatarURL:   "mxc://example.org/SEsfnsuifSDFSSEF",
			DisplayName: stringPtr("Alice Margatroid"),
			NewState:    MemberJoined,
			// TODO: Add Reason field.
			StateEventInfo: &StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomMember,
					},
					ID:               "$143273582443PhrSn:example.org",
					OriginServerTime: 1432735824653,
					RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
					Sender:           "@example:example.org",
					Unsigned: UnsignedData{
						Age: 1234,
					},
				},
				StateKey: "@alice:example.org",
			},
		},
	},
	{
		Name: "m.room.member 2",
		Code: `
			{
				"content": {
					"avatar_url": "mxc://example.org/SEsfnsuifSDFSSEF",
					"displayname": "Alice Margatroid",
					"membership": "invite",
					"reason": "Looking for support"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "@alice:example.org",
				"type": "m.room.member",
				"unsigned": {
					"age": 1234,
					"invite_room_state": [
						{
							"content": {
								"name": "Example Room"
							},
							"sender": "@bob:example.org",
							"state_key": "",
							"type": "m.room.name"
						},
						{
							"content": {
								"join_rule": "invite"
							},
							"sender": "@bob:example.org",
							"state_key": "",
							"type": "m.room.join_rules"
						}
					]
				}
			}
		`,
		Expected: RoomMemberEvent{
			AvatarURL:   "mxc://example.org/SEsfnsuifSDFSSEF",
			DisplayName: stringPtr("Alice Margatroid"),
			NewState:    MemberInvited,
			StateEventInfo: &StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomMember,
					},
					ID:               "$143273582443PhrSn:example.org",
					OriginServerTime: 1432735824653,
					RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
					Sender:           "@example:example.org",
					Unsigned: UnsignedData{
						Age: 1234,
						// TODO: Add field invite_room_state.
					},
				},
				StateKey: "@alice:example.org",
			},
		},
	},
	{
		Name: "m.room.member 3",
		Code: `
			{
				"content": {
					"avatar_url": "mxc://example.org/SEsfnsuifSDFSSEF",
					"displayname": "Alice Margatroid",
					"membership": "knock",
					"reason": "Looking for support"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "@alice:example.org",
				"type": "m.room.member",
				"unsigned": {
					"age": 1234,
					"knock_room_state": [
						{
							"content": {
								"name": "Example Room"
							},
							"sender": "@bob:example.org",
							"state_key": "",
							"type": "m.room.name"
						},
						{
							"content": {
								"join_rule": "knock"
							},
							"sender": "@bob:example.org",
							"state_key": "",
							"type": "m.room.join_rules"
						}
					]
				}
			}
		`,
		Expected: RoomMemberEvent{
			AvatarURL:   "mxc://example.org/SEsfnsuifSDFSSEF",
			DisplayName: stringPtr("Alice Margatroid"),
			NewState:    MemberKnock,
			// TODO: Add field Reason
			StateEventInfo: &StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomMember,
					},
					ID:               "$143273582443PhrSn:example.org",
					OriginServerTime: 1432735824653,
					RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
					Sender:           "@example:example.org",
				},
				StateKey: "@alice:example.org",
			},
		},
	},
	{
		Name: "m.room.member 4",
		Code: `
			{
				"content": {
					"avatar_url": "mxc://example.org/SEsfnsuifSDFSSEF",
					"displayname": "Alice Margatroid",
					"membership": "invite",
					"third_party_invite": {
						"display_name": "alice",
						"signed": {
							"mxid": "@alice:example.org",
							"signatures": {
								"magic.forest": {
									"ed25519:3": "fQpGIW1Snz+pwLZu6sTy2aHy/DYWWTspTJRPyNp0PKkymfIsNffysMl6ObMMFdIJhk6g6pwlIqZ54rxo8SLmAg"
								}
							},
							"token": "abc123"
						}
					}
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "@alice:example.org",
				"type": "m.room.member",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: RoomMemberEvent{
			AvatarURL:   "mxc://example.org/SEsfnsuifSDFSSEF",
			DisplayName: stringPtr("Alice Margatroid"),
			NewState:    MemberInvited,
			ThirdPartyInvite: struct {
				DisplayName string `json:"display_name"`
			}{
				DisplayName: "alice",
			},
			StateEventInfo: &StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomMember,
					},
					ID:               "$143273582443PhrSn:example.org",
					OriginServerTime: 1432735824653,
					RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
					Sender:           "@example:example.org",
					Unsigned: UnsignedData{
						Age: 1234,
					},
				},
				StateKey: "@alice:example.org",
			},
		},
	},
	{
		Name: "m.room.power_levels",
		Code: `
			{
				"content": {
					"ban": 50,
					"events": {
						"m.room.name": 100,
						"m.room.power_levels": 100
					},
					"events_default": 0,
					"invite": 50,
					"kick": 50,
					"notifications": {
						"room": 20
					},
					"redact": 50,
					"state_default": 50,
					"users": {
						"@example:localhost": 100
					},
					"users_default": 0
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.power_levels",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: RoomPowerLevelsEvent{
			BanRequirement: intPtr(50),
			Events: map[Type]int{
				TypeRoomName:        100,
				TypeRoomPowerLevels: 100,
			},
			EventRequirement:  0,
			InviteRequirement: intPtr(50),
			KickRequirement:   intPtr(50),
			Notifications: struct {
				Room *int `json:"room,omitempty"`
			}{
				Room: intPtr(20),
			},
			RedactRequirement: intPtr(50),
			StateRequirement:  50,
			UserLevel: map[matrix.UserID]int{
				"@example:localhost": 100,
			},
			UserDefault: 0,
			StateEventInfo: &StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomPowerLevels,
					},
					ID:               "$143273582443PhrSn:example.org",
					OriginServerTime: 1432735824653,
					RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
					Sender:           "@example:example.org",
					Unsigned: UnsignedData{
						Age: 1234,
					},
				},
				StateKey: "",
			},
		},
	},
	{
		Name: "m.room.redaction",
		Code: `
			{
				"content": {
					"reason": "Spamming"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"redacts": "$fukweghifu23:localhost",
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.room.redaction",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: RoomRedactionEvent{
			Reason: "Spamming",
			RoomEventInfo: &RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeRoomRedaction,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
		},
	},

	// TODO: Add other event types
}
