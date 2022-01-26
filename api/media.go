package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// MediaUpload uploads the provided file to the Matrix homeserver.
func (c *Client) MediaUpload(contentType string, filename string, body io.ReadCloser) (matrix.URL, error) {
	var resp struct {
		ContentURI matrix.URL `json:"content_uri"`
	}
	err := c.Request(
		"POST", c.Endpoints.MediaUpload(), &resp,
		httputil.WithToken(),
		httputil.WithHeader(map[string][]string{
			"Content-Type": {
				contentType,
			},
		}),
		httputil.WithQuery(map[string]string{
			"filename": filename,
		}),
		httputil.WithBody(body),
	)
	if err != nil {
		return "", fmt.Errorf("error uploading media: %w", err)
	}
	return resp.ContentURI, nil
}

// MediaDownloadURL returns the HTTP URL for the provided matrix URL.
// If allowRemote is false, the server will not attempt to fetch the media if it is deemed remote.
func (c *Client) MediaDownloadURL(matrixURL matrix.URL, allowRemote bool, filename string) (string, error) {
	parsed, err := url.Parse(string(matrixURL))
	if err != nil {
		return "", err
	}

	if parsed.Scheme != "mxc" {
		return string(matrixURL), nil
	}

	parsed.Path = strings.TrimPrefix(parsed.Path, "/")

	return c.FullRoute(c.Endpoints.MediaDownload(parsed.Host, parsed.Path, filename)) +
			"?allow_remote=" + strconv.FormatBool(allowRemote),
		nil
}

// MediaThumbnailMethod specifies the method the homeserver should crop the image in.
type MediaThumbnailMethod string

// The two types of valid values for MediaThumbnailMethod are scale and crop.
const (
	MediaThumbnailScale MediaThumbnailMethod = "scale"
	MediaThumbnailCrop  MediaThumbnailMethod = "crop"
)

// MediaThumbnailURL returns the HTTP URL for the provided matrix URL.
// If allowRemote is false, the server will not attempt to fetch the media if it is deemed remote.
// The provided width and height are treated as a guideline and the actual thumbnail may be a different size.
func (c *Client) MediaThumbnailURL(matrixURL matrix.URL, allowRemote bool,
	width int, height int, method MediaThumbnailMethod) (string, error) {
	parsed, err := url.Parse(string(matrixURL))
	if err != nil {
		return "", err
	}

	if parsed.Scheme != "mxc" {
		return string(matrixURL), nil
	}

	parsed.Path = strings.TrimPrefix(parsed.Path, "/")

	query := url.Values{
		"width":        {strconv.Itoa(width)},
		"height":       {strconv.Itoa(height)},
		"method":       {string(method)},
		"allow_remote": {strconv.FormatBool(allowRemote)},
	}

	return c.FullRoute(c.Endpoints.MediaThumbnail(parsed.Host, parsed.Path)) + "?" + query.Encode(), nil
}

// URLMetadata contains the basic OpenGraph metadata that the Matrix backend
// gives us using PreviewURL, as well as the raw JSON for the user to parse it
// further.
type URLMetadata struct {
	Title       string `json:"og:title,omitempty"`
	Type        string `json:"og:type,omitempty"`
	Description string `json:"og:description,omitempty"`
	URL         string `json:"og:url,omitempty"`

	Image       matrix.URL `json:"og:image,omitempty"`
	ImageType   string     `json:"og:image:type,omitempty"`
	ImageWidth  int        `json:"og:image:width,omitempty"`
	ImageHeight int        `json:"og:image:height,omitempty"`
	ImageSize   int        `json:"matrix:image:size,omitempty"`

	// Raw is the raw URL metadata received.
	Raw json.RawMessage `json:"-"`
}

// UnmarshalJSON parses b into u while keeping the raw JSON inside u.Raw.
func (u *URLMetadata) UnmarshalJSON(b []byte) error {
	// Copy the JSON.
	u.Raw = append(u.Raw[:0], b...)

	type urlMetadata URLMetadata
	return json.Unmarshal(b, (*urlMetadata)(u))
}

// MarshalJSON returns u.Raw if available, otherwise it marshals u.
func (u *URLMetadata) MarshalJSON() ([]byte, error) {
	if u.Raw != nil {
		return u.Raw, nil
	}

	type urlMetadata URLMetadata
	return json.Marshal((*urlMetadata)(u))
}

// PreviewURL requests the homeserver to generate a preview of the provided URL.
// It should be handled with care especially in an encrypted channel to prevent leaking URLs.
// ts is the preferred point in time to return a preview for. If it's zero value, the constraint is not passed on.
//
// The returned structure is a URLMetadata containing basic OpenGraph info, as well as the bundled
// JSON for further parsing.
func (c *Client) PreviewURL(url string, ts matrix.Timestamp) (*URLMetadata, error) {
	query := map[string]string{
		"url": url,
	}
	if ts != 0 {
		query["ts"] = strconv.Itoa(int(ts))
	}

	var resp *URLMetadata
	err := c.Request(
		"GET", c.Endpoints.MediaPreviewURL(), &resp,
		httputil.WithToken(), httputil.WithQuery(query),
	)
	if err != nil {
		return nil, fmt.Errorf("error previewing URL: %w", err)
	}

	return resp, nil
}

// MediaConfig is the configuration of the homeserver for media.
type MediaConfig struct {
	UploadSize int `json:"m.upload.size"`
}

// MediaConfig requests the media configuration of the server.
// Clients should follow the guide when using content repository endpoints.
func (c *Client) MediaConfig() (MediaConfig, error) {
	var resp MediaConfig
	err := c.Request(
		"GET", c.Endpoints.MediaConfig(), &resp,
		httputil.WithToken(),
	)
	if err != nil {
		return MediaConfig{}, fmt.Errorf("error fetching media config: %w", err)
	}
	return resp, nil
}
