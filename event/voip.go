package event

var (
	_ RoomEvent = CallInviteEvent{}
	_ RoomEvent = CallCandidatesEvent{}
	_ RoomEvent = CallAnswerEvent{}
	_ RoomEvent = CallHangupEvent{}
)

// TODO Maybe implement glare at some point.
// https://matrix.org/docs/spec/client_server/r0.6.1#glare

// CallInviteEvent is a message event where someone is inviting to establish a call.
//
// It has the type ID of `m.call.invite`.
type CallInviteEvent struct {
	RoomEventInfo

	CallID   string `json:"call_id"`
	Version  int    `json:"version"`  // Currently always 0.
	Lifetime int    `json:"lifetime"` // Milliseconds the offer is valid for.
	Offer    struct {
		Type string `json:"type"` // Must be "offer".
		SDP  string `json:"sdp"`  // Session Description Protocol
	} `json:"offer"`
}

// CallCandidatesEvent is a message event where additional ICE candidates are provided
// to foster communication.
//
// It has the type ID of `m.call.candidates`.
type CallCandidatesEvent struct {
	RoomEventInfo

	CallID     string `json:"call_id"`
	Version    int    `json:"version"` // Currently always 0.
	Candidates []struct {
		SDPMediaType      string `json:"sdpMid"`
		SDPMediaLineIndex int    `json:"sdpMLineIndex"`
		Candidate         string `json:"candidate"`
	} `json:"candidates"`
}

// CallAnswerEvent is a message event where a callee wishes to answer the call.
//
// It has the type ID of `m.call.answer`.
type CallAnswerEvent struct {
	RoomEventInfo

	CallID  string `json:"call_id"`
	Version int    `json:"int"`
	Answer  struct {
		Type string `json:"type"` // Must be "answer".
		SDP  string `json:"sdp"`  // Session Description Protocol
	} `json:"answer"`
}

// CallHangupEvent is a message event where the call is ended. This can be sent
// to hang up a call or to reject a call.
//
// It has the type ID of `m.call.hangup`.
type CallHangupEvent struct {
	RoomEventInfo

	CallID  string           `json:"call_id"`
	Version int              `json:"version"` // Currently always 0.
	Reason  CallHangupReason `json:"reason"`
}

// CallHangupReason is the reason we hung up.
type CallHangupReason string

// Possible reasons to hang up.
// List available at https://matrix.org/docs/spec/client_server/r0.6.1#m-call-hangup.
const (
	CallHangupNormal        CallHangupReason = ""
	CallHangupICEFailed     CallHangupReason = "ice_failed" // ICE negotiation failed.
	CallHangupInviteTimeout CallHangupReason = "invite_timeout"
)

// Type satisfies RoomEvent.
func (CallInviteEvent) Type() Type {
	return TypeCallInvite
}

// Type satisfies RoomEvent.
func (CallCandidatesEvent) Type() Type {
	return TypeCallCandidates
}

// Type satisfies RoomEvent.
func (CallAnswerEvent) Type() Type {
	return TypeCallAnswer
}

// Type satisfies RoomEvent.
func (CallHangupEvent) Type() Type {
	return TypeCallHangup
}

// SetRoomEventInfo sets the room event info.
func (c *CallInviteEvent) SetRoomEventInfo(i RoomEventInfo) {
	c.RoomEventInfo = i
}

// SetRoomEventInfo sets the room event info.
func (c *CallCandidatesEvent) SetRoomEventInfo(i RoomEventInfo) {
	c.RoomEventInfo = i
}

// SetRoomEventInfo sets the room event info.
func (c *CallAnswerEvent) SetRoomEventInfo(i RoomEventInfo) {
	c.RoomEventInfo = i
}

// SetRoomEventInfo sets the room event info.
func (c *CallHangupEvent) SetRoomEventInfo(i RoomEventInfo) {
	c.RoomEventInfo = i
}
