package event

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/chanbakjsd/gotrix/matrix"
)

func TestSpecExamples(t *testing.T) {
	for _, v := range specExamples {
		ev, err := RawEvent(v.Code).Parse()
		if err != nil {
			t.Errorf("error parsing %s: %v", v.Name, err)
			continue
		}

		v.Expected.Info().Raw = RawEvent(v.Code)
		if w, ok := ev.(*PresenceEvent); ok {
			// receiveTime is not deterministic. Remove the field for testing purpose after checking
			// it's filled.
			if w.receiveTime.IsZero() {
				t.Errorf("receiveTime is zero in PresenceEvent")
			}
			w.receiveTime = time.Time{}
		}

		if !reflect.DeepEqual(ev, v.Expected) {
			if bytes.Equal(v.Expected.Info().Raw, ev.Info().Raw) {
				// Redact raw if they're identical because it just spams output.
				v.Expected.Info().Raw = nil
				ev.Info().Raw = nil
			}

			t.Errorf("mismatch on parsing %s\nexpected: %#v\ngot: %#v", v.Name, v.Expected, ev)
			continue
		}
	}
}

var (
	boolTrue  = true
	boolFalse = false
)

func stringPtr(a string) *string {
	return &a
}

func intPtr(a int) *int {
	return &a
}

func urlPtr(a matrix.URL) *matrix.URL {
	return &a
}

func f64Ptr(a float64) *float64 {
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
		Expected: &RoomCanonicalAliasEvent{
			Alias: "#somewhere:localhost",
			AltAlias: []string{
				"#somewhere:example.org",
				"#myroom:example.com",
			},
			StateEventInfo: StateEventInfo{
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
		Expected: &RoomCreateEvent{
			Creator: "@example:example.org",
			// TODO: Add predecessor field
			Federated:   &boolTrue,
			RoomVersion: stringPtr("1"),
			StateEventInfo: StateEventInfo{
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
		Expected: &RoomJoinRulesEvent{
			JoinRule: JoinPublic,
			StateEventInfo: StateEventInfo{
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
		Expected: &RoomMemberEvent{
			UserID:      "@alice:example.org",
			AvatarURL:   "mxc://example.org/SEsfnsuifSDFSSEF",
			DisplayName: stringPtr("Alice Margatroid"),
			NewState:    MemberJoined,
			// TODO: Add Reason field.
			StateEventInfo: StateEventInfo{
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
		Expected: &RoomMemberEvent{
			UserID:      "@alice:example.org",
			AvatarURL:   "mxc://example.org/SEsfnsuifSDFSSEF",
			DisplayName: stringPtr("Alice Margatroid"),
			NewState:    MemberInvited,
			StateEventInfo: StateEventInfo{
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
		Expected: &RoomMemberEvent{
			UserID:      "@alice:example.org",
			AvatarURL:   "mxc://example.org/SEsfnsuifSDFSSEF",
			DisplayName: stringPtr("Alice Margatroid"),
			NewState:    MemberKnock,
			// TODO: Add field Reason
			StateEventInfo: StateEventInfo{
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
		Expected: &RoomMemberEvent{
			UserID:      "@alice:example.org",
			AvatarURL:   "mxc://example.org/SEsfnsuifSDFSSEF",
			DisplayName: stringPtr("Alice Margatroid"),
			NewState:    MemberInvited,
			ThirdPartyInvite: struct {
				DisplayName string `json:"display_name"`
			}{
				DisplayName: "alice",
			},
			StateEventInfo: StateEventInfo{
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
		Expected: &RoomPowerLevelsEvent{
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
			StateEventInfo: StateEventInfo{
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
		Expected: &RoomRedactionEvent{
			Redacts: "$fukweghifu23:localhost",
			Reason:  "Spamming",
			RoomEventInfo: RoomEventInfo{
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
	{
		Name: "m.room.message audio",
		Code: `
			{
				"content": {
					"body": "Bee Gees - Stayin' Alive",
					"info": {
						"duration": 2140786,
						"mimetype": "audio/mpeg",
						"size": 1563685
					},
					"msgtype": "m.audio",
					"url": "mxc://example.org/ffed755USFFxlgbQYZGtryd"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.room.message",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomMessageEvent{
			Body: "Bee Gees - Stayin' Alive",
			AdditionalInfo: json.RawMessage(`{
						"duration": 2140786,
						"mimetype": "audio/mpeg",
						"size": 1563685
					}`),
			MessageType: RoomMessageAudio,
			URL:         "mxc://example.org/ffed755USFFxlgbQYZGtryd",
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeRoomMessage,
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
	{
		Name: "m.room.message emote",
		Code: `
		{
			"content": {
				"body": "thinks this is an example emote",
				"format": "org.matrix.custom.html",
				"formatted_body": "thinks <b>this</b> is an example emote",
				"msgtype": "m.emote"
			},
			"event_id": "$143273582443PhrSn:example.org",
			"origin_server_ts": 1432735824653,
			"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
			"sender": "@example:example.org",
			"type": "m.room.message",
			"unsigned": {
				"age": 1234
			}
		}
		`,
		Expected: &RoomMessageEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeRoomMessage,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			Body:          "thinks this is an example emote",
			Format:        FormatHTML,
			FormattedBody: "thinks <b>this</b> is an example emote",
			MessageType:   RoomMessageEmote,
		},
	},
	{
		Name: "m.room.message file",
		Code: `
			{
				"content": {
					"body": "something-important.doc",
					"filename": "something-important.doc",
					"info": {
						"mimetype": "application/msword",
						"size": 46144
					},
					"msgtype": "m.file",
					"url": "mxc://example.org/FHyPlCeYUSFFxlgbQYZmoEoe"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.room.message",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomMessageEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeRoomMessage,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			Body: "something-important.doc",
			// TODO: Add field Filename
			AdditionalInfo: json.RawMessage(`{
						"mimetype": "application/msword",
						"size": 46144
					}`),
			MessageType: RoomMessageFile,
			URL:         "mxc://example.org/FHyPlCeYUSFFxlgbQYZmoEoe",
		},
	},
	{
		Name: "m.room.message image",
		Code: `
			{
				"content": {
					"body": "filename.jpg",
					"info": {
						"h": 398,
						"mimetype": "image/jpeg",
						"size": 31037,
						"w": 394
					},
					"msgtype": "m.image",
					"url": "mxc://example.org/JWEIFJgwEIhweiWJE"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.room.message",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomMessageEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeRoomMessage,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			Body: "filename.jpg",
			AdditionalInfo: json.RawMessage(`{
						"h": 398,
						"mimetype": "image/jpeg",
						"size": 31037,
						"w": 394
					}`),
			MessageType: RoomMessageImage,
			URL:         "mxc://example.org/JWEIFJgwEIhweiWJE",
		},
	},
	{
		Name: "m.room.message location",
		Code: `
			{
				"content": {
					"body": "Big Ben, London, UK",
					"geo_uri": "geo:51.5008,0.1247",
					"info": {
						"thumbnail_info": {
							"h": 300,
							"mimetype": "image/jpeg",
							"size": 46144,
							"w": 300
						},
						"thumbnail_url": "mxc://example.org/FHyPlCeYUSFFxlgbQYZmoEoe"
					},
					"msgtype": "m.location"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.room.message",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomMessageEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeRoomMessage,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			Body:   "Big Ben, London, UK",
			GeoURI: "geo:51.5008,0.1247",
			AdditionalInfo: json.RawMessage(`{
						"thumbnail_info": {
							"h": 300,
							"mimetype": "image/jpeg",
							"size": 46144,
							"w": 300
						},
						"thumbnail_url": "mxc://example.org/FHyPlCeYUSFFxlgbQYZmoEoe"
					}`),
			MessageType: RoomMessageLocation,
		},
	},
	{
		Name: "m.room.message notice",
		Code: `
			{
				"content": {
					"body": "This is an example notice",
					"format": "org.matrix.custom.html",
					"formatted_body": "This is an <strong>example</strong> notice",
					"msgtype": "m.notice"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.room.message",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomMessageEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeRoomMessage,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			Body:          "This is an example notice",
			Format:        FormatHTML,
			FormattedBody: "This is an <strong>example</strong> notice",
			MessageType:   RoomMessageNotice,
		},
	},
	{
		Name: "m.room.message server notice",
		Code: `
			{
				"content": {
					"admin_contact": "mailto:server.admin@example.org",
					"body": "Human-readable message to explain the notice",
					"limit_type": "monthly_active_user",
					"msgtype": "m.server_notice",
					"server_notice_type": "m.server_notice.usage_limit_reached"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.room.message",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomMessageEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeRoomMessage,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			// TODO: Add admin_contact, limit_type, server_notice_type field
			Body:        "Human-readable message to explain the notice",
			MessageType: "m.server_notice",
		},
	},
	{
		Name: "m.room.message text",
		Code: `
			{
				"content": {
					"body": "This is an example text message",
					"format": "org.matrix.custom.html",
					"formatted_body": "<b>This is an example text message</b>",
					"msgtype": "m.text"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.room.message",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomMessageEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeRoomMessage,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			Body:          "This is an example text message",
			Format:        FormatHTML,
			FormattedBody: "<b>This is an example text message</b>",
			MessageType:   RoomMessageText,
		},
	},
	{
		Name: "m.room.message video",
		Code: `
			{
				"content": {
					"body": "Gangnam Style",
					"info": {
						"duration": 2140786,
						"h": 320,
						"mimetype": "video/mp4",
						"size": 1563685,
						"thumbnail_info": {
							"h": 300,
							"mimetype": "image/jpeg",
							"size": 46144,
							"w": 300
						},
						"thumbnail_url": "mxc://example.org/FHyPlCeYUSFFxlgbQYZmoEoe",
						"w": 480
					},
					"msgtype": "m.video",
					"url": "mxc://example.org/a526eYUSFFxlgbQYZmo442"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.room.message",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomMessageEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeRoomMessage,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			Body: "Gangnam Style",
			AdditionalInfo: json.RawMessage(`{
						"duration": 2140786,
						"h": 320,
						"mimetype": "video/mp4",
						"size": 1563685,
						"thumbnail_info": {
							"h": 300,
							"mimetype": "image/jpeg",
							"size": 46144,
							"w": 300
						},
						"thumbnail_url": "mxc://example.org/FHyPlCeYUSFFxlgbQYZmoEoe",
						"w": 480
					}`),
			MessageType: RoomMessageVideo,
			URL:         "mxc://example.org/a526eYUSFFxlgbQYZmo442",
		},
	},
	{
		Name: "m.room.name",
		Code: `
			{
				"content": {
					"name": "The room name"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.name",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomNameEvent{
			StateEventInfo: StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomName,
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
			Name: "The room name",
		},
	},
	{
		Name: "m.room.topic",
		Code: `
			{
				"content": {
					"topic": "A room topic"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.topic",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomTopicEvent{
			StateEventInfo: StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomTopic,
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
			Topic: "A room topic",
		},
	},
	{
		Name: "m.room.avatar",
		Code: `
			{
				"content": {
					"info": {
						"h": 398,
						"mimetype": "image/jpeg",
						"size": 31037,
						"w": 394
					},
					"url": "mxc://example.org/JWEIFJgwEIhweiWJE"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.avatar",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomAvatarEvent{
			StateEventInfo: StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomAvatar,
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
			Image: ImageInfo{
				FileInfo: FileInfo{
					MimeType: "image/jpeg",
					Size:     31037,
				},
				Height: 398,
				Width:  394,
			},
			URL: "mxc://example.org/JWEIFJgwEIhweiWJE",
		},
	},
	{
		Name: "m.room.pinned_events",
		Code: `
			{
				"content": {
					"pinned": [
						"$someevent:example.org"
					]
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.pinned_events",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomPinnedEvent{
			StateEventInfo: StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomPinned,
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
			Pinned: []matrix.EventID{
				"$someevent:example.org",
			},
		},
	},
	{
		Name: "m.direct",
		Code: `
			{
				"content": {
					"@bob:example.com": [
						"!abcdefgh:example.com",
						"!hgfedcba:example.com"
					]
				},
				"type": "m.direct"
			}
		`,
		Expected: &DirectEvent{
			EventInfo: EventInfo{
				Type: TypeDirect,
			},
			Rooms: map[matrix.UserID][]matrix.RoomID{
				"@bob:example.com": {
					"!abcdefgh:example.com",
					"!hgfedcba:example.com",
				},
			},
		},
	},
	{
		Name: "m.call.invite",
		Code: `
			{
				"content": {
					"call_id": "12345",
					"lifetime": 60000,
					"offer": {
						"sdp": "v=0\r\no=- 6584580628695956864 2 IN IP4 127.0.0.1[...]",
						"type": "offer"
					},
					"version": 0
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.call.invite",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &CallInviteEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeCallInvite,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			CallID:   "12345",
			Lifetime: 60000,
			Offer: struct {
				Type string `json:"type"`
				SDP  string `json:"sdp"`
			}{
				Type: "offer",
				SDP:  "v=0\r\no=- 6584580628695956864 2 IN IP4 127.0.0.1[...]",
			},
			Version: 0,
		},
	},
	{
		Name: "m.call.candidates",
		Code: `
			{
				"content": {
					"call_id": "12345",
					"candidates": [
						{
							"candidate": "candidate:863018703 1 udp 2122260223 10.9.64.156 43670 typ host generation 0",
							"sdpMLineIndex": 0,
							"sdpMid": "audio"
						}
					],
					"version": 0
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.call.candidates",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &CallCandidatesEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeCallCandidates,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			CallID: "12345",
			Candidates: []struct {
				SDPMediaType      string `json:"sdpMid"`
				SDPMediaLineIndex int    `json:"sdpMLineIndex"`
				Candidate         string `json:"candidate"`
			}{
				{
					Candidate:         "candidate:863018703 1 udp 2122260223 10.9.64.156 43670 typ host generation 0",
					SDPMediaLineIndex: 0,
					SDPMediaType:      "audio",
				},
			},
			Version: 0,
		},
	},
	{
		Name: "m.call.answer",
		Code: `
			{
				"content": {
					"answer": {
						"sdp": "v=0\r\no=- 6584580628695956864 2 IN IP4 127.0.0.1[...]",
						"type": "answer"
					},
					"call_id": "12345",
					"lifetime": 60000,
					"version": 0
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.call.answer",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &CallAnswerEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeCallAnswer,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			Answer: struct {
				Type string `json:"type"`
				SDP  string `json:"sdp"`
			}{
				Type: "answer",
				SDP:  "v=0\r\no=- 6584580628695956864 2 IN IP4 127.0.0.1[...]",
			},
			CallID: "12345",
			// Lifetime is deliberately not included because it was not specified in spec.
			Version: 0,
		},
	},
	{
		Name: "m.call.hangup",
		Code: `
			{
				"content": {
					"call_id": "12345",
					"version": 0
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"type": "m.call.hangup",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &CallHangupEvent{
			RoomEventInfo: RoomEventInfo{
				EventInfo: EventInfo{
					Type: TypeCallHangup,
				},
				ID:               "$143273582443PhrSn:example.org",
				OriginServerTime: 1432735824653,
				RoomID:           "!jEsUZKDJdhlrceRyVU:example.org",
				Sender:           "@example:example.org",
				Unsigned: UnsignedData{
					Age: 1234,
				},
			},
			CallID:  "12345",
			Version: 0,
		},
	},
	{
		Name: "m.typing",
		Code: `
			{
			  "content": {
			    "user_ids": [
			      "@alice:matrix.org",
			      "@bob:example.com"
			    ]
			  },
			  "room_id": "!jEsUZKDJdhlrceRyVU:example.org",
			  "type": "m.typing"
			}
		`,
		Expected: &TypingEvent{
			EventInfo: EventInfo{
				Type: TypeTyping,
			},
			UserID: []matrix.UserID{
				"@alice:matrix.org",
				"@bob:example.com",
			},
			RoomID: "!jEsUZKDJdhlrceRyVU:example.org",
		},
	},
	{
		Name: "m.receipt",
		Code: `
			{
				"content": {
					"$1435641916114394fHBLK:matrix.org": {
						"m.read": {
							"@rikj:jki.re": {
								"ts": 1436451550453
							}
						}
					}
				},
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"type": "m.receipt"
			}
		`,
		Expected: &ReceiptEvent{
			EventInfo: EventInfo{
				Type: TypeReceipt,
			},
			Events: map[matrix.EventID]Receipt{
				"$1435641916114394fHBLK:matrix.org": {
					Read: map[matrix.UserID]struct {
						Timestamp int `json:"ts"`
					}{
						"@rikj:jki.re": {
							Timestamp: 1436451550453,
						},
					},
				},
			},
			RoomID: "!jEsUZKDJdhlrceRyVU:example.org",
		},
	},
	{
		Name: "m.presence",
		Code: `
			{
				"content": {
					"avatar_url": "mxc://localhost/wefuiwegh8742w",
					"currently_active": false,
					"last_active_ago": 2478593,
					"presence": "online",
					"status_msg": "Making cupcakes"
				},
				"sender": "@example:localhost",
				"type": "m.presence"
			}
		`,
		Expected: &PresenceEvent{
			EventInfo: EventInfo{
				Type: TypePresence,
			},
			AvatarURL:       urlPtr("mxc://localhost/wefuiwegh8742w"),
			CurrentlyActive: &boolFalse,
			LastActiveAgo:   intPtr(2478593),
			Presence:        matrix.PresenceOnline,
			Status:          stringPtr("Making cupcakes"),
			User:            "@example:localhost",
		},
	},
	{
		Name: "m.room.history_visibility",
		Code: `
			{
				"content": {
					"history_visibility": "shared"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.history_visibility",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomHistoryVisibilityEvent{
			StateEventInfo: StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomHistoryVisibility,
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
			Visibility: VisibilityShared,
		},
	},
	{
		Name: "m.room.guest_access",
		Code: `
			{
				"content": {
					"guest_access": "can_join"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.guest_access",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomGuestAccessEvent{
			StateEventInfo: StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomGuestAccess,
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
			GuestAccess: GuestAccessCanJoin,
		},
	},
	{
		Name: "m.tag",
		Code: `
			{
				"content": {
					"tags": {
						"u.work": {
							"order": 0.9
						}
					}
				},
				"type": "m.tag"
			}
		`,
		Expected: &TagEvent{
			EventInfo: EventInfo{
				Type: TypeTag,
			},
			Tags: map[matrix.TagName]matrix.Tag{
				"u.work": {
					Order: f64Ptr(0.9),
				},
			},
		},
	},
	{
		Name: "m.room.tombstone",
		Code: `
			{
				"content": {
					"body": "This room has been replaced",
					"replacement_room": "!newroom:example.org"
				},
				"event_id": "$143273582443PhrSn:example.org",
				"origin_server_ts": 1432735824653,
				"room_id": "!jEsUZKDJdhlrceRyVU:example.org",
				"sender": "@example:example.org",
				"state_key": "",
				"type": "m.room.tombstone",
				"unsigned": {
					"age": 1234
				}
			}
		`,
		Expected: &RoomTombstoneEvent{
			StateEventInfo: StateEventInfo{
				RoomEventInfo: RoomEventInfo{
					EventInfo: EventInfo{
						Type: TypeRoomTombstone,
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
			Message:         "This room has been replaced",
			ReplacementRoom: "!newroom:example.org",
		},
	},

	// TODO: Add other event types
}
