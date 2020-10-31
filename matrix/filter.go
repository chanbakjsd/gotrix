package matrix

// Filter represents a filter that can be uploaded to/downloaded from
// the homeserver.
//
// Servers MAY still send data that has been excluded by the filter.
// The filter only tells the server what is safe to not include.
type Filter struct {
	// List of event fields that should be included.
	EventFields []string `json:"event_fields,omitempty"`
	// The format to use. This implementation recommends the
	// use of "client" (default).
	//
	// Don't include if you don't know what you're doing.
	EventFormat string `json:"event_format,omitempty"`
	// List of presence updates to include.
	Presence EventFilter `json:"presence,omitempty"`
	// List of user account data updates to include.
	// This does not affect data associated with room.
	AccountData EventFilter `json:"account_data,omitempty"`
	// Filter to be applied to room data.
	Room RoomFilter `json:"room,omitempty"`
}

// EventFilter represents a filter that filters events that should
// be sent to the client.
type EventFilter struct {
	// Maximum number of events to return.
	Limit int `json:"limit,omitempty"`
	// List of senders to include. All if not provided.
	IncludedSenders []string `json:"senders,omitempty"`
	// List of event types to include. All if not provided.
	// '*' can be used as a wildcard.
	IncludedTypes []string `json:"types,omitempty"`
	// List of senders to exclude. Overrides IncludedSenders.
	ExcludedSenders []string `json:"not_senders,omitempty"`
	// List of event types to exclude. Overrides IncludedTypes.
	// '*' can be used as a wildcard.
	ExcludedTypes []string `json:"not_types,omitempty"`
}

// RoomFilter represents a filter that filters room data.
type RoomFilter struct {
	// Rooms to include. All if not provided.
	IncludedRooms []string `json:"rooms,omitempty"`
	// Rooms to exclude. Overrides IncludedRoom.
	ExcludedRooms []string `json:"not_rooms,omitempty"`
	// Ephemeral is the subfilter applied to events that are
	// not persistent (added to history) like typing.
	Ephemeral RoomEventFilter `json:"ephemeral,omitempty"`
	// The client will continue to listen to events from rooms
	// that the user has left if this is set to true.
	// Defaults to false.
	IncludeLeave bool `json:"include_leave,omitempty"`
	// State is the subfilter applied to state events.
	State StateFilter `json:"state,omitempty"`
	// Timeline is the subfilter applied to events that are
	// persistent (added to history) like messages.
	Timeline RoomEventFilter `json:"timeline"`
	// AccountData is the subfilter applied to per user account
	// data.
	AccountData RoomEventFilter `json:"account_data,omitempty"`
}

// StateFilter represents a filter for state events.
type StateFilter struct {
	// Limit is the maximum number of events to return.
	Limit int `json:"limit,omitempty"`
	// List of senders to include. All if omitted.
	IncludedSenders []string `json:"senders,omitempty"`
	// List of types to include. All if omitted.
	IncludedTypes []string `json:"types"`
	// List of rooms to include. All if omitted.
	IncludeRooms []string `json:"rooms,omitempty"`
	// List of senders to exclude. Overrides IncludedSenders.
	ExcludedSenders []string `json:"not_senders,omitempty"`
	// List of types to exclude. Overrides IncludedTypes.
	ExcludedTypes []string `json:"not_types,omitempty"`
	// List of rooms to exclude. Overrides IncludedRooms.
	ExcludedRooms []string `json:"not_rooms,omitempty"`
	// Enable lazy loading members. If it's true, it'll only send
	// member info that are mentioned in events.
	// Other member data should be queried through the API if this
	// is true.
	LazyLoadMembers bool `json:"lazy_load_members,omitempty"`
	// The server does not send member info that it thinks the client
	// already knows by default. The server will include it instead if
	// this is set to true.
	IncludeRedundantMembers bool `json:"include_redundant_members,omitempty"`
	// Include only events with a `url` key in its content if `true`.
	// Include only events without a `url` key in its content if `false`.
	// `url` is not used to filter otherwise.
	ContainsURL *bool `json:"contains_url,omitempty"`
}

// RoomEventFilter represents a filter that filters room events.
type RoomEventFilter struct {
	// Limit is the maximum number of events to return.
	Limit int `json:"limit,omitempty"`
	// List of senders to include. All if omitted.
	IncludedSenders []string `json:"senders,omitempty"`
	// List of types to include. All if omitted.
	IncludedTypes []string `json:"types"`
	// List of rooms to include. All if omitted.
	IncludeRooms []string `json:"rooms,omitempty"`
	// List of senders to exclude. Overrides IncludedSenders.
	ExcludedSenders []string `json:"not_senders,omitempty"`
	// List of types to exclude. Overrides IncludedTypes.
	ExcludedTypes []string `json:"not_types,omitempty"`
	// List of rooms to exclude. Overrides IncludedRooms.
	ExcludedRooms []string `json:"not_rooms,omitempty"`
	// Enable lazy loading members. If it's true, it'll only send
	// member info that are mentioned in events.
	// Other member data should be queried through the API if this
	// is true.
	LazyLoadMembers bool `json:"lazy_load_members,omitempty"`
	// The server does not send member info that it thinks the client
	// already knows by default. The server will include it instead if
	// this is set to true.
	IncludeRedundantMembers bool `json:"include_redundant_members,omitempty"`
	// Include only events with a `url` key in its content if `true`.
	// Include only events without a `url` key in its content if `false`.
	// `url` is not used to filter otherwise.
	ContainsURL *bool `json:"contains_url,omitempty"`
}
