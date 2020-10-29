package gomatrix

import (
	"github.com/chanbakjsd/gomatrix/api"
)

// Client is an instance of a higher level client.
type Client struct {
	*api.Client
}
