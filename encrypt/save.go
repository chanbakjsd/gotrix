package encrypt

import (
	"encoding/json"
	"fmt"

	"github.com/chanbakjsd/gotrix/encrypt/olm"
)

type storage struct {
	Account string `json:"acc"`
}

// RestoreSession restores the provided data to a Session. The key provided
// should be equivalent to the key passed to (*Session).Save when generating
// data.
func RestoreSession(c *Client, opts Options, data []byte, key string) (*Session, error) {
	var store storage
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	account, err := olm.NewAccountFromPickle(key, store.Account)
	if err != nil {
		return nil, err
	}
	s := &Session{
		Client:  c,
		opts:    opts,
		account: account,
	}
	if err := s.setupHooks(); err != nil {
		return nil, err
	}
	return s, nil
}

// Save stores the private keys associated with the session and encrypt them
// using the provided key.
func (s *Session) Save(key string) ([]byte, error) {
	pickledAccount, err := s.account.Pickle(key)
	if err != nil {
		return nil, fmt.Errorf("error pickling account: %w", err)
	}
	return json.Marshal(storage{
		Account: pickledAccount,
	})
}
