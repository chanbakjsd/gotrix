package api

import (
	"fmt"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// User is one user as returned by UserSearch.
type User struct {
	ID          matrix.UserID `json:"user_id"`
	DisplayName *string       `json:"display_name"`
	AvatarURL   *matrix.URL   `json:"avatar_url"`
}

// UserSearch searches for users that match the keyword. The default limit is 10 if the zero value is provided.
// 'limited' is true when the result is being limited by the provided limit.
func (c *Client) UserSearch(keyword string, limit int) (result []User, limited bool, err error) {
	req := struct {
		SearchTerm string `json:"search_term"`
		Limit      int    `json:"limit"`
	}{
		SearchTerm: keyword,
		Limit:      limit,
	}
	var resp struct {
		Results []User `json:"results"`
		Limited bool   `json:"limited"`
	}

	err = c.Request(
		"POST", EndpointUserDirectorySearch, &resp,
		httputil.WithToken(), httputil.WithJSONBody(req),
	)
	if err != nil {
		return nil, false, fmt.Errorf("error searching for user: %w", err)
	}

	return resp.Results, resp.Limited, nil
}

// User returns the combined info of the provided user.
func (c *Client) User(userID matrix.UserID) (User, error) {
	resp := User{
		ID: userID,
	}
	err := c.Request(
		"GET", EndpointProfile(userID), &resp,
	)
	if err != nil {
		return User{}, fmt.Errorf("error fetching user info: %w", err)
	}
	return resp, nil
}

// DisplayName returns the display name of the provided user.
func (c *Client) DisplayName(userID matrix.UserID) (*string, error) {
	var resp struct {
		DisplayName *string `json:"displayname,omitempty"`
	}
	err := c.Request(
		"GET", EndpointProfileDisplayName(userID), &resp,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching display name: %w", err)
	}
	return resp.DisplayName, nil
}

// DisplayNameSet sets the display name of the provided user.
func (c *Client) DisplayNameSet(displayName string) error {
	req := map[string]string{
		"displayname": displayName,
	}

	err := c.Request(
		"PUT", EndpointProfileDisplayName(c.UserID), nil,
		httputil.WithToken(), httputil.WithJSONBody(req),
	)
	if err != nil {
		return fmt.Errorf("error setting display name: %w", err)
	}
	return nil
}

// AvatarURL returns the avatar URL of the provided user.
func (c *Client) AvatarURL(userID matrix.UserID) (*matrix.URL, error) {
	var resp struct {
		AvatarURL *matrix.URL `json:"avatar_url,omitempty"`
	}
	err := c.Request(
		"GET", EndpointProfileAvatarURL(userID), &resp,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching avatar URL: %w", err)
	}
	return resp.AvatarURL, nil
}

// AvatarURLSet sets the avatar URL of the provided user.
func (c *Client) AvatarURLSet(avatarURL matrix.URL) error {
	req := map[string]interface{}{
		"avatar_url": avatarURL,
	}

	err := c.Request(
		"PUT", EndpointProfileAvatarURL(c.UserID), nil,
		httputil.WithToken(), httputil.WithJSONBody(req),
	)
	if err != nil {
		return fmt.Errorf("error setting avatar URL: %w", err)
	}
	return nil
}
