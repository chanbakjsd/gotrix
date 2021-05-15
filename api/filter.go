package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/event"
)

// FilterAdd uploads the provided filter to the homeserver and returns its
// assigned ID.
func (c *Client) FilterAdd(filterToUpload event.GlobalFilter) (string, error) {
	var resp struct {
		FilterID string `json:"filter_id"`
	}
	err := c.Request(
		"POST", EndpointFilter(c.UserID), &resp,
		httputil.WithToken(), httputil.WithJSONBody(filterToUpload),
	)
	if err != nil {
		return "", fmt.Errorf("error adding filter: %w", err)
	}
	return resp.FilterID, nil
}

// Filter downloads the requested filter from the homeserver.
func (c *Client) Filter(filterID string) (*event.GlobalFilter, error) {
	resp := &event.GlobalFilter{}
	err := c.Request(
		"GET", EndpointFilterGet(c.UserID, filterID), resp,
		httputil.WithToken(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting filter: %w", err)
	}
	return resp, nil
}
