package api

import (
	"errors"
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

var ErrInvalidTagName = errors.New("tag name must not exceed 255 bytes")

// Tags retrives the tags for a room.
func (c *Client) Tags(roomID matrix.RoomID) (map[matrix.TagName]matrix.Tag, error) {
	var resp struct {
		Tags map[matrix.TagName]matrix.Tag `json:"tags"`
	}
	err := c.Request(
		"GET", EndpointTags(c.UserID, roomID), &resp,
		httputil.WithToken(),
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching tags: %w", err)
	}
	return resp.Tags, nil
}

// TagAdd adds a room tag.
func (c *Client) TagAdd(roomID matrix.RoomID, name matrix.TagName, tagData matrix.Tag) error {
	if len([]byte(name)) > 255 {
		return ErrInvalidTagName
	}
	err := c.Request(
		"PUT", EndpointTag(c.UserID, roomID, name), nil,
		httputil.WithToken(), httputil.WithJSONBody(tagData),
	)
	if err != nil {
		return fmt.Errorf("error adding tag: %w", err)
	}
	return nil
}

// TagDelete removes a room tag.
func (c *Client) TagDelete(roomID matrix.RoomID, name matrix.TagName) error {
	if len([]byte(name)) > 255 {
		return ErrInvalidTagName
	}
	err := c.Request(
		"DELETE", EndpointTag(c.UserID, roomID, name), nil,
		httputil.WithToken(),
	)
	if err != nil {
		return fmt.Errorf("error deleting tag: %w", err)
	}
	return nil
}
