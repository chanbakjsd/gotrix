// Package encrypt provides end-to-end encryption support for the Matrix API.
package encrypt

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/chanbakjsd/gotrix/api"
	"github.com/chanbakjsd/gotrix/matrix"
)

// SignatureAlgorithm is the type that describes the names of signature algorithms.
type SignatureAlgorithm string

const (
	Ed25519    SignatureAlgorithm = "ed25519"
	Curve25519 SignatureAlgorithm = "curve25519"
)

// SignedAlgorithm returns a signed variant of the signature algorithm. Calling this method multiple
// times will not repeat the same signed prefix.
func (a SignatureAlgorithm) SignedAlgorithm() SignatureAlgorithm {
	if !strings.HasPrefix(string(a), "signed_") {
		return "signed_" + a
	}
	return a
}

// Signature describes a digital cryptography signature of unknown algorithm.
type Signature string

// SyncResponse contains the extended parts of api.SyncResponse, holding encrypt-specific fields.
// Use ParseSyncResponse to extract these fields.
//
// Quoting from the documentation:
//
// This module adds an optional device_lists property to the /sync response, as specified below. The
// server need only populate this property for an incremental /sync (i.e., one where the since
// parameter was specified). The client is expected to use /keys/query or /keys/changes for the
// equivalent functionality after an initial sync, as documented in Tracking the device list for a
// user.
type SyncResponse struct {
	// DeviceLists is the information on e2e device updates. Note: only present on an incremental
	// sync.
	DeviceLists struct {
		// Changed is the list of users who have updated their device identity or cross-signing
		// keys, or who now share an encrypted room with the client since the previous sync
		// response.
		Changed []matrix.UserID `json:"changed"`
		// Left is the list of users with whom we do not share any encrypted rooms anymore since the
		// previous sync response.
		Left []matrix.UserID `json:"left"`
	} `json:"device_lists,omitempty"`
	// DeviceOneTimeKeysCount is the number of unclaimed one-time keys currently held on the server
	// for this device.
	DeviceOneTimeKeysCount map[Algorithm]int `json:"device_one_time_keys_count,omitempty"`
}

// ParseSyncResponse parses the encrypt.SyncResponse from the given API SyncResponse.
func ParseSyncResponse(resp *api.SyncResponse) (*SyncResponse, error) {
	if resp.Raw() == nil {
		return nil, errors.New("SyncResponse missing raw JSON")
	}

	var encResp SyncResponse

	if err := json.Unmarshal(resp.Raw(), &encResp); err != nil {
		return nil, err
	}

	return &encResp, nil
}
