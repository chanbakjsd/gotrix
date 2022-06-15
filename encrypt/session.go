package encrypt

import (
	"encoding/json"
	"fmt"

	ejson "github.com/chanbakjsd/gotrix/encrypt/json"
	"github.com/chanbakjsd/gotrix/encrypt/olm"
)

// Options is a list of configurable options for the encrypt session.
type Options struct {
	// AutoDecrypt will decrypt all encryption events and forward it as regular
	// events.
	AutoDecrypt bool
}

// Session is an instance of a device in the end-to-end encryption setup.
type Session struct {
	*Client
	opts Options

	account *olm.Account
}

// NewSession adds hooks related to encryption to provide the specified client
// with end-to-end encryption support and returns a client wrapper.
//
// To resume a session, use the RestoreSession function.
//
// This function assumes that the client is already authenticated and makes
// necessary API calls to establish itself as a new device.
func NewSession(c *Client, opts Options) (*Session, error) {
	account, err := olm.NewAccount()
	if err != nil {
		return nil, err
	}
	s := &Session{
		Client:  c,
		opts:    opts,
		account: account,
	}
	if err := s.uploadIdentityKeys(); err != nil {
		return nil, err
	}
	if err := s.setupHooks(); err != nil {
		return nil, err
	}
	return s, nil
}

// SignData signs the provided data by marshalling it in the canonical form and
// passing it to libolm.
func (s *Session) SignData(data interface{}) (Signature, error) {
	json, err := ejson.CanonicalMarshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling data to sign: %w", err)
	}
	signature, err := s.account.Sign(string(json))
	if err != nil {
		return "", fmt.Errorf("error signing marshalled JSON: %w", err)
	}
	return Signature(signature), nil
}

// uploadIdentityKeys uploads the identity keys of the session's account.
func (s *Session) uploadIdentityKeys() error {
	idKeysJSON, err := s.account.IdentityKeys()
	if err != nil {
		return err
	}

	var idKeys map[SignatureAlgorithm]Key
	if err := json.Unmarshal([]byte(idKeysJSON), &idKeys); err != nil {
		return fmt.Errorf("error unmarshalling generated identity keys: %w", err)
	}
	keyToUpload := make(map[DeviceAlgorithmKey]Key, len(idKeys))
	for k, v := range idKeys {
		keyToUpload[NewDeviceAlgorithmKey(s.DeviceID, k)] = v
	}

	deviceKeys := DeviceKeys{
		Algorithms: []Algorithm{AlgorithmOlm, AlgorithmMegOlm},
		DeviceID:   s.DeviceID,
		Keys:       keyToUpload,
		UserID:     s.UserID,
	}
	signature, err := s.SignData(deviceKeys)
	if err != nil {
		return fmt.Errorf("error signing generated identity keys: %w", err)
	}
	deviceKeys.Signatures = UserSignatures{
		s.UserID: {
			NewDeviceAlgorithmKey(s.DeviceID, Ed25519): signature,
		},
	}
	if _, err = s.KeyUpload(deviceKeys, OneTimeKeys{}); err != nil {
		return fmt.Errorf("error uploading identity keys: %w", err)
	}
	return nil
}
