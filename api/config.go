package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// ClientConfig retrieves the client config previously stored with ClientConfigSet.
// 'v' is a pointer that is directly passed into json.Unmarshal.
func (c *Client) ClientConfig(userID matrix.UserID, configType string, v interface{}) error {
	err := c.Request(
		"PUT", EndpointAccountDataGlobal(userID, configType), v,
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
func (c *Client) ClientConfigSet(userID matrix.UserID, configType string, config interface{}) error {
	err := c.Request(
		"PUT", EndpointAccountDataGlobal(userID, configType), nil,
		httputil.WithToken(), httputil.WithJSONBody(config),
	)
	if err != nil {
		return fmt.Errorf("error saving client config: %w", err)
	}
	return nil
}

// ClientConfigRoom retrieves the client config previously stored with ClientConfigRoomSet.
// 'v' is a pointer that is directly passed into json.Unmarshal.
func (c *Client) ClientConfigRoom(userID matrix.UserID, roomID matrix.RoomID, configType string, v interface{}) error {
	err := c.Request(
		"PUT", EndpointAccountDataRoom(userID, roomID, configType), v,
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
func (c *Client) ClientConfigRoomSet(userID matrix.UserID, roomID matrix.RoomID, configType string,
	config interface{}) error {
	err := c.Request(
		"PUT", EndpointAccountDataRoom(userID, roomID, configType), nil,
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
func (c *Client) IgnoredUsers(userID matrix.UserID) ([]matrix.UserID, error) {
	var resp ignoredUsers
	err := c.ClientConfig(userID, "m.ignored_user_list", &resp)
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
func (c *Client) IgnoredUsersSet(userID matrix.UserID, newList []matrix.UserID) error {
	req := ignoredUsers{
		IgnoredUsers: make(map[matrix.UserID]struct{}),
	}
	for _, v := range newList {
		req.IgnoredUsers[v] = struct{}{}
	}

	err := c.ClientConfigSet(userID, "m.ignored_user_list", req)
	if err != nil {
		return fmt.Errorf("error setting ignored users: %w", err)
	}
	return nil
}
