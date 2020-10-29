package api

import (
	"github.com/chanbakjsd/gomatrix/matrix"
)

// Error represents an API error as returned by the Matrix server.
type Error struct {
	// ResponseCode is the status code of the response and is not
	// part of the actual response.
	ResponseCode int `json:"-"`

	// Code and Message should be included in every API error.
	Code    matrix.ErrorCode `json:"errcode"`
	Message string           `json:"error"`

	// SoftLogout is included in invalid token errors.
	// If it's true, the client should just log back in.
	// If it's false, the client should purge all its cache before
	// logging back in.
	SoftLogout bool `json:"soft_logout"`

	// RetryAfterMillisecond is included in rate limit errors.
	RetryAfterMillisecond int `json:"retry_after_ms"`
}

// Error makes API Error implement the `error` interface.
func (e Error) Error() string {
	return e.Message
}

// HTTPError represents an error while decoding response.
// It contains the status code and the actual error.
type HTTPError struct {
	Code            int
	UnderlyingError error
}

// Error makes HTTPError implement the `error` interface.
func (h HTTPError) Error() string {
	return h.UnderlyingError.Error()
}

// Unwrap allows the underlying error to be exposed.
func (h HTTPError) Unwrap() error {
	return h.UnderlyingError
}
