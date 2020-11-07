package gomatrix

import (
	"net/url"
	"reflect"

	"github.com/chanbakjsd/gomatrix/api"
	"github.com/chanbakjsd/gomatrix/api/httputil"
	"github.com/chanbakjsd/gomatrix/event"
)

// Client is an instance of a higher level client.
type Client struct {
	*api.Client
	Filter  event.GlobalFilter
	Handler Handler

	nextRetryTime int
	shouldClose   bool
	closeDone     chan struct{}
}

// New creates a client with the provided host URL and the default HTTP client.
// It assumes https if the scheme is not provided.
func New(homeServerHost string) (*Client, error) {
	return NewWithClient(httputil.NewClient(), homeServerHost)
}

// NewWithClient creates a client with the provided host URL and the provided client.
// It assumes https if the scheme is not provided.
func NewWithClient(httpClient httputil.Client, serverName string) (*Client, error) {
	parsed, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme == "" {
		parsed.Scheme = "https"
	}

	apiClient := &api.Client{
		Client: httpClient,
	}
	apiClient.HomeServer = parsed.Host
	apiClient.HomeServerScheme = parsed.Scheme
	return &Client{
		Client: apiClient,
		Filter: DefaultFilter,
		Handler: &defaultHandler{
			handlers: make(map[event.Type][]reflect.Value),
		},
	}, err
}
