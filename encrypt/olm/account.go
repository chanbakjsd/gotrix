package olm

// #include <stdlib.h>
// #include "olm/olm.h"
import "C"

import (
	"crypto/rand"
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// Account is an account that is associated with a device owned by a user.
type Account struct {
	acc *C.OlmAccount
}

// NewRawAccount creates a new uninitialized account.
func NewRawAccount() *Account {
	mem := C.malloc(accountSize)
	cAcc := C.olm_account(mem)
	acc := Account{
		acc: cAcc,
	}
	runtime.SetFinalizer(&acc, func(_ *Account) {
		C.free(mem)
	})

	return &acc
}

// NewAccount creates a new account initialized by random bytes from crypto/rand.
func NewAccount() (*Account, error) {
	a := NewRawAccount()
	randomLen := C.olm_create_account_random_length(a.acc)
	randomBytes := make([]byte, randomLen)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("olm: error fetching random bytes: %w", err)
	}

	cBytes := C.CBytes(randomBytes)
	defer C.free(cBytes)

	ret := C.olm_create_account(a.acc, cBytes, randomLen)
	if ret == errValue {
		return nil, a.LastError()
	}
	return a, nil
}

// NewAccountFromPickle creates an account from the provided key and pickle.
func NewAccountFromPickle(key, pickle string) (*Account, error) {
	a := NewRawAccount()
	cKey := unsafe.Pointer(C.CString(key))
	cPickle := unsafe.Pointer(C.CString(pickle))
	defer C.free(cKey)
	defer C.free(cPickle)

	ret := C.olm_unpickle_account(a.acc, cKey, C.size_t(len(key)), cPickle, C.size_t(len(pickle)))
	if ret == errValue {
		return nil, a.LastError()
	}
	return a, nil
}

// LastError returns the last error Account returned.
func (a *Account) LastError() error {
	str := C.GoString(C.olm_account_last_error(a.acc))
	if str == noErr {
		return nil
	}
	return errors.New("olm: account error: " + str)
}

// Clear clears the backing memory of Account.
func (a *Account) Clear() error {
	ret := C.olm_clear_account(a.acc)
	if ret == errValue {
		return a.LastError()
	}
	return nil
}

// Pickle stores the account as a base64 string encrypted by key.
func (a *Account) Pickle(key string) (string, error) {
	length := C.olm_pickle_account_length(a.acc)
	pickle := C.malloc(length)
	cKey := C.CString(key)
	defer C.free(pickle)
	defer C.free(unsafe.Pointer(cKey))

	ret := C.olm_pickle_account(a.acc, unsafe.Pointer(cKey), C.size_t(len(key)), pickle, length)
	if ret == errValue {
		return "", a.LastError()
	}
	return C.GoStringN((*C.char)(pickle), C.int(ret)), nil
}

// IdentityKeys returns the public part of the identity keys for an account.
func (a *Account) IdentityKeys() (string, error) {
	length := C.olm_account_identity_keys_length(a.acc)
	idKeys := C.malloc(length)
	defer C.free(idKeys)

	ret := C.olm_account_identity_keys(a.acc, idKeys, length)
	if ret == errValue {
		return "", a.LastError()
	}
	return C.GoString((*C.char)(idKeys)), nil
}

// Sign signs the provided message with the ED25519 key of the account.
func (a *Account) Sign(message string) (string, error) {
	length := C.olm_account_signature_length(a.acc)
	cMsg := unsafe.Pointer(C.CString(message))
	signature := C.malloc(length)
	defer C.free(cMsg)
	defer C.free(signature)

	ret := C.olm_account_sign(a.acc, cMsg, C.size_t(len(message)), signature, length)
	if ret == errValue {
		return "", a.LastError()
	}
	return C.GoString((*C.char)(signature)), nil
}

// MaxOneTimeKeys returns the maximum number of one time keys.
func (a *Account) MaxOneTimeKeys() int {
	return int(C.olm_account_max_number_of_one_time_keys(a.acc))
}

// OneTimeKeys generates a new set of one time keys and returns unpublished one time keys
// in the form of a JSON-formatted object containing the 'curve25519' property.
//
// The 'curve25519' property contains an object that maps key ID to base64 encoded Curve25519 keys.
func (a *Account) OneTimeKeys(count int) (string, error) {
	cCount := C.ulong(count)
	randomLen := C.olm_account_generate_one_time_keys_random_length(a.acc, cCount)
	randomBytes := make([]byte, randomLen)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("olm: error fetching random bytes: %w", err)
	}

	cBytes := C.CBytes(randomBytes)
	defer C.free(cBytes)

	ret := C.olm_account_generate_one_time_keys(a.acc, cCount, cBytes, randomLen)
	if ret == errValue {
		return "", a.LastError()
	}

	length := C.olm_account_one_time_keys_length(a.acc)
	oneTimeKeys := C.malloc(length)
	defer C.free(oneTimeKeys)

	ret = C.olm_account_one_time_keys(a.acc, oneTimeKeys, length)
	if ret == errValue {
		return "", a.LastError()
	}
	return C.GoString((*C.char)(oneTimeKeys)), nil
}

// MarkKeysAsPublished marks the current set of one time keys as published.
func (a *Account) MarkKeysAsPublished() error {
	ret := C.olm_account_mark_keys_as_published(a.acc)
	if ret == errValue {
		return a.LastError()
	}
	return nil
}

// RemoveOneTimeKeys removes the one-time keys associated with the provided Session.
func (a *Account) RemoveOneTimeKeys(s *Session) error {
	ret := C.olm_remove_one_time_keys(a.acc, s.sess)
	if ret == errValue {
		return a.LastError()
	}
	return nil
}
