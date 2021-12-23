package api

import (
	"encoding/json"
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/matrix"
)

// ClientConfig retrieves the client config previously stored with ClientConfigSet.
// 'v' is a pointer that is directly passed into json.Unmarshal to unmarshal the content of the config.
func (c *Client) ClientConfig(configType string, v interface{}) error {
	err := c.Request(
		"GET", EndpointAccountDataGlobal(c.UserID, configType), v,
		httputil.WithToken(),
	)
	if err != nil {
		return fmt.Errorf("error getting client config: %w", err)
	}
	return nil
}

// ClientConfigSet sets the client config to the provided value.
// 'configType' should be namespaced Java-style (com.example.someData) to prevent clashes.
// 'config' is JSON-encoded before sent to the server.
//
// The provided config will be provided as an event with type 'configType' and content 'config' in AccountData
// of SyncResponse.
func (c *Client) ClientConfigSet(configType string, config interface{}) error {
	err := c.Request(
		"PUT", EndpointAccountDataGlobal(c.UserID, configType), nil,
		httputil.WithToken(), httputil.WithJSONBody(config),
	)
	if err != nil {
		return fmt.Errorf("error saving client config: %w", err)
	}
	return nil
}

// ClientConfigRoom retrieves the client config previously stored with ClientConfigRoomSet.
// 'v' is a pointer that is directly passed into json.Unmarshal to unmarshal the content of the config.
func (c *Client) ClientConfigRoom(roomID matrix.RoomID, configType string, v interface{}) error {
	err := c.Request(
		"GET", EndpointAccountDataRoom(c.UserID, roomID, configType), v,
		httputil.WithToken(),
	)
	if err != nil {
		return fmt.Errorf("error getting client config: %w", err)
	}
	return nil
}

// ClientConfigRoomSet sets the client config to the provided value. It is equivalent to ClientConfigSet except it is
// scoped to a Matrix room.
// 'configType' should be namespaced Java-style (com.example.someData) to prevent clashes.
// 'config' is JSON-encoded before sent to the server.
//
// The provided config will be provided as an event with type 'configType' and content 'config' in AccountData
// of rooms in SyncResponse.
func (c *Client) ClientConfigRoomSet(roomID matrix.RoomID, configType string,
	config interface{}) error {
	err := c.Request(
		"PUT", EndpointAccountDataRoom(c.UserID, roomID, configType), nil,
		httputil.WithToken(), httputil.WithJSONBody(config),
	)
	if err != nil {
		return fmt.Errorf("error saving client config: %w", err)
	}
	return nil
}

type ignoredUsers struct {
	IgnoredUsers map[matrix.UserID]struct{} `json:"ignored_users"`
}

// IgnoredUsers returns the list of users configured to be ignored.
func (c *Client) IgnoredUsers() ([]matrix.UserID, error) {
	var resp ignoredUsers
	err := c.ClientConfig("m.ignored_user_list", &resp)
	if err != nil {
		return nil, fmt.Errorf("error getting ignored users: %w", err)
	}

	list := make([]matrix.UserID, 0, len(resp.IgnoredUsers))
	for k := range resp.IgnoredUsers {
		list = append(list, k)
	}
	return list, nil
}

// IgnoredUsersSet sets the list of users configured to be ignored.
func (c *Client) IgnoredUsersSet(newList []matrix.UserID) error {
	req := ignoredUsers{
		IgnoredUsers: make(map[matrix.UserID]struct{}),
	}
	for _, v := range newList {
		req.IgnoredUsers[v] = struct{}{}
	}

	err := c.ClientConfigSet("m.ignored_user_list", req)
	if err != nil {
		return fmt.Errorf("error setting ignored users: %w", err)
	}
	return nil
}

// DMRooms fetches the list of DM rooms as saved in 'm.direct'.
func (c *Client) DMRooms() (*event.DirectEvent, error) {
	var resp event.RawEvent
	err := c.ClientConfig("m.direct", &resp)
	if err != nil {
		return nil, fmt.Errorf("error fetching DM room list: %w", err)
	}

	ev, err := event.Parse(resp)
	if err != nil {
		return nil, fmt.Errorf("error parsing DM room list: %w", err)
	}

	directEvent, ok := ev.(*event.DirectEvent)
	if !ok {
		return nil, fmt.Errorf("error parsing DM room list: got %T instead of m.direct", ev)
	}
	return directEvent, nil
}

// DMRoomsSet updates the DM rooms saved in 'm.direct'.
func (c *Client) DMRoomsSet(newRooms *event.DirectEvent) error {
	raw, err := json.Marshal(newRooms)
	if err != nil {
		return fmt.Errorf("error encoding DM rooms: %w", err)
	}
	err = c.ClientConfigSet("m.direct", raw)
	if err != nil {
		return fmt.Errorf("error setting DM rooms: %w", err)
	}
	return nil
}
