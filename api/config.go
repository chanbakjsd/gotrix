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
