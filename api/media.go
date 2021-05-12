package api

import (
	"errors"
	"io"
	"net/url"
	"strconv"

	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

var (
	// ErrNoUploadPerm means that the user does not have permission to upload the content.
	// The file type may not be permitted by the server. The user may have reached a quota for uploaded content.
	ErrNoUploadPerm = errors.New("no upload permission")
	// ErrFileTooLarge means that the file is too large for the server.
	ErrFileTooLarge = errors.New("file is too large")
)

// MediaUpload uploads the provided file to the Matrix homeserver.
//
// It implements the `POST _matrix/media/r0/upload` endpoint.
func (c *Client) MediaUpload(contentType string, filename string, body io.ReadCloser) (matrix.URL, error) {
	var resp struct {
		ContentURI matrix.URL `json:"content_uri"`
	}
	err := c.Request(
		"POST", "_matrix/media/r0/upload", &resp,
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

	return resp.ContentURI, matrix.MapAPIError(err, matrix.ErrorMap{
		matrix.CodeForbidden: ErrNoUploadPerm,
		matrix.CodeTooLarge:  ErrFileTooLarge,
	})
}

// MediaDownloadURL returns the HTTP URL for the provided matrix URL.
// If allowRemote is false, the server will not attempt to fetch the media if it is deemed remote.
//
// It implements the `POST _matrix/media/r0/download/{serverName}/{mediaId}/{fileName}` endpoint.
func (c *Client) MediaDownloadURL(matrixURL matrix.URL, allowRemote bool, filename string) (string, error) {
	parsed, err := url.Parse(string(matrixURL))
	if err != nil {
		return "", err
	}

	if parsed.Scheme != "mxc" {
		return string(matrixURL), nil
	}

	return c.HomeServerScheme + "://" + c.HomeServer + "/_matrix/media/r0/download/" +
			url.PathEscape(parsed.Host) + "/" + url.PathEscape(parsed.Path) + "/" + url.PathEscape(filename) +
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
//
// It implements the `POST _matrix/media/r0/thumbnail/{serverName}/{mediaId}` endpoint.
func (c *Client) MediaThumbnailURL(matrixURL matrix.URL, allowRemote bool,
	width int, height int, method MediaThumbnailMethod) (string, error) {
	parsed, err := url.Parse(string(matrixURL))
	if err != nil {
		return "", err
	}

	if parsed.Scheme != "mxc" {
		return string(matrixURL), nil
	}

	return c.HomeServerScheme + "://" + c.HomeServer + "/_matrix/media/r0/thumbnail/" +
			url.PathEscape(parsed.Host) + "/" + url.PathEscape(parsed.Path) + "?" +
			"width=" + strconv.Itoa(width) + "&height=" + strconv.Itoa(height) +
			"&method=" + url.QueryEscape(string(method)) + "&allow_remote=" + strconv.FormatBool(allowRemote),
		nil
}

// PreviewURL requests the homeserver to generate a preview of the provided URL.
// It should be handled with care especially in an encrypted channel to prevent leaking URLs.
// ts is the preferred point in time to return a preview for. If it's zero value, the constraint is not passed on.
//
// The returned map is a map of OpenGraph info. og:image will be an MXC URI to the image instead if available.
//
// It implements the `GET _matrix/media/r0/preview_url` endpoint.
func (c *Client) PreviewURL(url string, ts matrix.Timestamp) (map[string]interface{}, error) {
	query := map[string]string{
		"url": url,
	}
	if ts != 0 {
		query["ts"] = strconv.Itoa(int(ts))
	}

	var resp map[string]interface{}
	err := c.Request(
		"GET", "_matrix/media/r0/preview_url", &resp,
		httputil.WithToken(), httputil.WithQuery(query),
	)

	return resp, err
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
		"GET", "_matrix/media/r0/config", &resp,
		httputil.WithToken(),
	)
	return resp, err
}
