package gotrix

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/chanbakjsd/gotrix/api"
	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/event"
	"github.com/chanbakjsd/gotrix/state"
)

// Client is an instance of a higher level client.
type Client struct {
	*api.Client

	SyncOpts SyncOptions
	Handler  Handler
	State    State

	ctx        context.Context
	next       string
	cancelFunc func()
	closeDone  chan struct{}
}

// New creates a client with the provided host URL and the default HTTP client.
// It assumes https if the scheme is not provided.
func New(homeServerHost string) (*Client, error) {
	return NewWithClient(httputil.NewClient(), homeServerHost)
}

// NewWithClient creates a client with the provided host URL and the provided client.
// It assumes https if the scheme is not provided.
func NewWithClient(httpClient httputil.Client, serverName string) (*Client, error) {
	if !strings.Contains(serverName, "://") {
		// First is protocol while second is port.
		serverName = "https://" + serverName
	}
	parsed, err := url.Parse(serverName)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %s: %w", serverName, err)
	}

	apiClient := &api.Client{
		Client:    httpClient,
		Endpoints: api.Endpoints{Version: "r0"},
	}
	apiClient.HomeServer = parsed.Host
	apiClient.HomeServerScheme = parsed.Scheme

	if vClient, err := apiClient.WithLatestVersion(); err == nil {
		apiClient = vClient
	}

	return &Client{
		Client:   apiClient,
		SyncOpts: DefaultSyncOptions,
		Handler: &defaultHandler{
			handlers: make(map[event.Type][]reflect.Value),
		},
		State: state.NewDefault(),
	}, nil
}

// WithContext creates a copy of the client that uses the provided context.
func (c Client) WithContext(ctx context.Context) *Client {
	c.Client = c.Client.WithContext(ctx)
	c.ctx = ctx
	return &c
}
