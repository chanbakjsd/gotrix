package gomatrix

import (
	"github.com/chanbakjsd/gomatrix/api"
	"github.com/chanbakjsd/gomatrix/matrix"
)

// LoginPassword authenticates the client using the provided username and password.
func (c *Client) LoginPassword(username, password string) error {
	return c.Login(api.LoginArg{
		Type: matrix.LoginPassword,
		Identifier: matrix.Identifier{
			Type: matrix.IdentifierUser,
			User: username,
		},
		Password: password,
	})
}

// LoginToken authenticates the client using the provided token.
func (c *Client) LoginToken(token string) error {
	return c.Login(api.LoginArg{
		Type:  matrix.LoginToken,
		Token: token,
	})
}
