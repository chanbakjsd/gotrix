package api

import (
	"net/url"

	"github.com/chanbakjsd/gomatrix/api/httputil"
	"github.com/chanbakjsd/gomatrix/event"
)

// FilterAdd uploads the provided filter to the homeserver and returns its
// assigned ID.
//
// It implements the `POST _matrix/client/r0/user/{userId}/filter` endpoint.
func (c *Client) FilterAdd(filterToUpload event.GlobalFilter) (string, error) {
	var resp struct {
		FilterID string `json:"filter_id"`
	}
	err := c.Request(
		"POST", "_matrix/client/r0/user/"+url.PathEscape(string(c.UserID))+"/filter", &resp,
		httputil.WithToken(),
		httputil.WithBody(filterToUpload),
	)
	return resp.FilterID, err
}

// Filter downloads the requested filter from the homeserver.
//
// It implements the `GET _matrix/client/r0/user/{userId}/filter/{filterId}`
// endpoint.
func (c *Client) Filter(filterID string) (*event.GlobalFilter, error) {
	resp := &event.GlobalFilter{}
	err := c.Request(
		"GET", "_matrix/client/r0/user/"+url.PathEscape(string(c.UserID))+"/filter/"+filterID, resp,
		httputil.WithToken(),
	)

	return resp, err
}
