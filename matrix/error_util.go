package matrix

import (
	"errors"
)

// StatusCode takes in an error and return the HTTP status code associated with it.
//
// If it's not a HTTPError, it returns -1 instead.
func StatusCode(e error) int {
	var err HTTPError
	if errors.As(e, &err) {
		return err.Code
	}
	return -1
}

// ErrCode takes in an error and return the API error code associated with it.
//
// If it's not an APIError, it returns an empty string instead.
func ErrCode(e error) ErrorCode {
	var err APIError
	if errors.As(e, &err) {
		return err.Code
	}
	return ""
}

// ErrorMap is a shorthand for map[ErrorCode]error.
type ErrorMap map[ErrorCode]error

// MapAPIError is a helper function that maps API errors to its concrete error types as
// provided by the user.
//
// Unmatched errors are returned as-is.
func MapAPIError(e error, m ErrorMap) error {
	var err APIError
	if !errors.As(e, &err) {
		return e
	}
	if x, ok := m[err.Code]; ok {
		return x
	}
	return e
}
