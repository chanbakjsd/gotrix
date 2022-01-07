package encrypt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chanbakjsd/gotrix"
	"github.com/chanbakjsd/gotrix/api/httputil"
	"github.com/chanbakjsd/gotrix/matrix"
)

// Client wraps around an api.Client to provide additional encryption endpoint
// calls.
type Client struct {
	*gotrix.Client
	Endpoints Endpoints
}

// NewClient creates a new Client wrapping the provided gotrix.Client.
func NewClient(c *gotrix.Client) (*Client, error) {
	return &Client{
		Client: c,
		Endpoints: Endpoints{
			Endpoints: c.Endpoints,
		},
	}, nil
}

// KeyChanges is returned by encrypt.Client.KeyChanges.
type KeyChanges struct {
	Changed []matrix.UserID `json:"changed"`
	Left    []matrix.UserID `json:"left"`
}

// KeyChanges gets a list of users who have updated their device identity keys
// since a previous sync token.
func (c *Client) KeyChanges(from, to string) (*KeyChanges, error) {
	var resp KeyChanges

	err := c.Request(
		"GET", c.Endpoints.KeysChanges(), &resp,
		httputil.WithToken(),
		httputil.WithQuery(map[string]string{
			"from": from,
			"to":   to,
		}),
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// KeyClaims contains the keys to be claimed. It maps from device ID to
// algorithm name.
type KeyClaims map[matrix.UserID]string

// KeyClaimResponse is returned by encrypt.Client.KeyClaim.
type KeyClaimResponse struct {
	OneTimeKeys OneTimeUserDeviceKeys `json:"one_time_keys"`
	// TODO: Failures, but the exact type is not documented.
}

// OneTimeUserDeviceKeys is the map of one-time keys for the queried devices.
// It maps the user ID given to Client.KeyClaim to a map that maps the device
// ID to the user-device keys.
type OneTimeUserDeviceKeys map[matrix.UserID]map[matrix.DeviceID]OneTimeKeys

// OneTimeKeys maps an object with key "<algorithm>:<keyID>" to the user-device
// key object.
type OneTimeKeys map[string]json.RawMessage

// NewOneTimeKeys creates a map of one-time keys from the given two key maps.
func NewOneTimeKeys(
	unsigned map[Algorithm]map[matrix.DeviceID]Key,
	signed map[Algorithm]map[matrix.DeviceID]OneTimeSignedKey,
) OneTimeKeys {
	otk := OneTimeKeys{}

	for algo, devices := range unsigned {
		for device, key := range devices {
			b, err := json.Marshal(key)
			if err != nil {
				panic("cannot marshal unsigned key: " + err.Error())
			}
			otk[string(algo)+":"+string(device)] = b
		}
	}

	for algo, devices := range signed {
		for device, key := range devices {
			b, err := json.Marshal(key)
			if err != nil {
				panic("cannot marshal signed key: " + err.Error())
			}
			otk[string(algo)+":"+string(device)] = b
		}
	}

	return otk
}

// Key gets the unsigned key from the given signature algorithm type and key
// ID.
func (ks OneTimeKeys) Key(algorithm SignatureAlgorithm, keyID string) (Key, error) {
	raw := ks.RawKey(algorithm.SignedAlgorithm(), keyID)
	if raw == nil {
		return "", fmt.Errorf("no key with ID %q found", keyID)
	}

	var key Key
	if err := json.Unmarshal(raw, &key); err != nil {
		return "", err
	}

	return key, nil
}

// SignedKey gets the signed key from the given signature algorithm type and
// key ID. The user does not have to call SignedAlgorithm on the given
// algorithm.
func (ks OneTimeKeys) SignedKey(algorithm SignatureAlgorithm, keyID string) (OneTimeSignedKey, error) {
	raw := ks.RawKey(algorithm.SignedAlgorithm(), keyID)
	if raw == nil {
		return OneTimeSignedKey{}, fmt.Errorf("no key with ID %q found", keyID)
	}

	var otk OneTimeSignedKey

	if err := json.Unmarshal(raw, &otk); err != nil {
		return OneTimeSignedKey{}, err
	}

	return otk, nil
}

// RawKey gets the key value as raw JSON from the given signature algorithm
// type and key ID.
func (ks OneTimeKeys) RawKey(algorithm SignatureAlgorithm, keyID string) json.RawMessage {
	return ks[string(algorithm)+":"+keyID]
}

// OneTimeSignedKey describes the value type of a "signed" key.
type OneTimeSignedKey struct {
	Key        Key
	Signatures UserSignatures
}

// DeviceAlgorithmKey combines both the signature algorithm and device ID to
// form a colon-delimited key string used in the API. It has the format
// "<algorithm>:<deviceID>".
type DeviceAlgorithmKey string

// NewDeviceAlgorithmKey creates a newly-formatted DeviceAlgorithmKey.
func NewDeviceAlgorithmKey(deviceID matrix.DeviceID, algorithm SignatureAlgorithm) DeviceAlgorithmKey {
	return DeviceAlgorithmKey(string(algorithm) + ":" + string(deviceID))
}

// UserSignatures maps from user ID to a map of "<algorithm>:<deviceID>" to a
// signature.
type UserSignatures map[matrix.UserID]map[DeviceAlgorithmKey]Signature

// Signature gets the signature for the given user, device and algorithm.
func (s UserSignatures) Signature(uID matrix.UserID, deviceID matrix.DeviceID, algorithm SignatureAlgorithm) Signature {
	return s[uID][NewDeviceAlgorithmKey(deviceID, algorithm)]
}

// KeyClaim claims one-time keys for use in pre-key messages. The timeout has
// an accuracy of 1ms and 10 seconds is the recommended default.
func (c *Client) KeyClaim(userClaims map[matrix.UserID]KeyClaims, timeout time.Duration) (*KeyClaimResponse, error) {
	req := struct {
		OneTimeKeys map[matrix.UserID]KeyClaims `json:"one_time_keys"`
		Timeout     int                         `json:"timeout,omitempty"`
	}{
		OneTimeKeys: userClaims,
		Timeout:     int(timeout / time.Millisecond),
	}

	var resp KeyClaimResponse

	err := c.Request(
		"GET", c.Endpoints.KeysClaim(), &resp,
		httputil.WithToken(),
		httputil.WithJSONBody(req),
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// KeyQueryResponse is returned from encrypt.Client.KeyQuery.
type KeyQueryResponse struct {
	DeviceKeys      map[matrix.UserID]map[matrix.DeviceID]DeviceKeys `json:"device_keys"`
	MasterKeys      map[matrix.UserID]MasterKeys                     `json:"master_keys"`
	SelfSigningKeys map[matrix.UserID]SelfSigningKeys                `json:"self_signing_keys"`
	UserSigningKeys map[matrix.UserID]UserSigningKeys                `json:"user_signing_keys"`
	// TODO: failures
}

// DeviceKeys contains information on a queried device.
type DeviceKeys struct {
	Algorithms []Algorithm                `json:"algorithms"`
	DeviceID   matrix.DeviceID            `json:"device_id"`
	Keys       map[DeviceAlgorithmKey]Key `json:"keys"`
	Signatures UserSignatures             `json:"signatures"`
	UserID     matrix.UserID              `json:"user_id"`
	Unsigned   json.RawMessage            `json:"unsigned,omitempty"` // not in KeyUpload
}

// MasterKeys contains information on the master cross-signing keys.
type MasterKeys struct {
	Keys   map[string]string `json:"keys"`
	Usage  []string          `json:"usage"` // []string{"master"}
	UserID matrix.UserID     `json:"user_id"`
}

// SelfSigningKeys contains information on the self-signing keys.
type SelfSigningKeys struct {
	Keys       map[string]string `json:"keys"`
	Signatures UserSignatures    `json:"signatures"`
	Usage      []string          `json:"usage"`
	UserID     matrix.UserID     `json:"user_id"`
}

// UserSigningKeys contains information on the user-signing key of the user
// making the request, if they queried their own device information
type UserSigningKeys struct {
	Keys       map[string]string `json:"keys"`
	Signatures UserSignatures    `json:"signatures"`
	Usage      []string          `json:"usage"`
	UserID     matrix.UserID     `json:"user_id"`
}

// KeyQuery returns the current devices and identity keys for the given users.
// The timeout has an accuracy of 1ms and has a recommended default of 10
// seconds.
//
// If the client is fetching keys as a result of a device update received in a
// sync request, token should be the 'since' token of that sync request, or any
// later sync token. This allows the server to ensure its response contains the
// keys advertised by the notification in that sync.
func (c *Client) KeyQuery(deviceKeys map[matrix.UserID][]matrix.DeviceID, timeout time.Duration, token string) (*KeyQueryResponse, error) {
	req := struct {
		DeviceKeys map[matrix.UserID][]matrix.DeviceID `json:"device_keys"`
		Timeout    int                                 `json:"timeout,omitempty"`
		Token      string                              `json:"token,omitempty"`
	}{
		DeviceKeys: deviceKeys,
		Timeout:    int(timeout / time.Millisecond),
		Token:      token,
	}

	var resp KeyQueryResponse

	err := c.Request(
		"GET", c.Endpoints.KeysQuery(), &resp,
		httputil.WithToken(),
		httputil.WithJSONBody(req),
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// KeyUpload publishes end-to-end encryption keys for the device. It returns the new count of
// one-time keys for the device given in deviceKeys.
func (c *Client) KeyUpload(deviceKeys DeviceKeys, oneTimeKeys OneTimeKeys) (count map[Algorithm]int, err error) {
	req := struct {
		DeviceKeys  DeviceKeys  `json:"device_keys"`
		OneTimeKeys OneTimeKeys `json:"one_time_keys"`
	}{
		DeviceKeys:  deviceKeys,
		OneTimeKeys: oneTimeKeys,
	}

	var resp struct {
		OneTimeKeyCounts map[Algorithm]int `json:"one_time_key_counts"`
	}

	err = c.Request(
		"GET", c.Endpoints.KeysUpload(), &resp,
		httputil.WithToken(),
		httputil.WithJSONBody(req),
	)
	if err != nil {
		return nil, err
	}

	return resp.OneTimeKeyCounts, nil
}
