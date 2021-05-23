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

// Session is a struct that tracks a group of people that share the same keys with each other.
type Session struct {
	sess *C.OlmSession
}

// NewRawSession creates a new uninitialized session.
func NewRawSession() *Session {
	mem := C.malloc(sessionSize)
	cSess := C.olm_session(mem)
	s := Session{
		sess: cSess,
	}
	runtime.SetFinalizer(&s, func(_ *Session) {
		C.free(mem)
	})

	return &s
}

// NewOutboundSession generates a new outbound session.
func NewOutboundSession(a *Account, theirIdentityKey string, theirOneTimeKey string) (*Session, error) {
	s := NewRawSession()
	randomLen := C.olm_create_outbound_session_random_length(s.sess)
	randomBytes := make([]byte, randomLen)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("olm: error fetching random bytes: %w", err)
	}

	cBytes := C.CBytes(randomBytes)
	defer C.free(cBytes)

	cIDKey := unsafe.Pointer(C.CString(theirIdentityKey))
	cOTKey := unsafe.Pointer(C.CString(theirOneTimeKey))
	defer C.free(cIDKey)
	defer C.free(cOTKey)

	ret := C.olm_create_outbound_session(
		s.sess, a.acc, cIDKey, C.size_t(len(theirIdentityKey)), cOTKey, C.size_t(len(theirOneTimeKey)),
		cBytes, randomLen,
	)
	if ret == errValue {
		return nil, s.LastError()
	}
	return s, nil
}

// NewInboundSession generates a new inbound session from the local Account, the other party's identity key
// and the message send by the other party.
func NewInboundSession(a *Account, identityKey string, otkMsg string) (*Session, error) {
	s := NewRawSession()

	cIDKey := unsafe.Pointer(C.CString(identityKey))
	cOTKMsg := unsafe.Pointer(C.CString(otkMsg))
	defer C.free(cIDKey)
	defer C.free(cOTKMsg)

	ret := C.olm_create_inbound_session_from(
		s.sess, a.acc, cIDKey, C.size_t(len(identityKey)), cOTKMsg, C.size_t(len(otkMsg)),
	)
	if ret == errValue {
		return nil, s.LastError()
	}
	return s, nil
}

// NewSessionFromPickle creates a session from the provided key and pickle.
func NewSessionFromPickle(key, pickle string) (*Session, error) {
	s := NewRawSession()
	cKey := unsafe.Pointer(C.CString(key))
	cPickle := unsafe.Pointer(C.CString(pickle))
	defer C.free(cKey)
	defer C.free(cPickle)

	ret := C.olm_unpickle_session(s.sess, cKey, C.size_t(len(key)), cPickle, C.size_t(len(pickle)))
	if ret == errValue {
		return nil, s.LastError()
	}
	return s, nil
}

// LastError returns the last error Session returned.
func (s *Session) LastError() error {
	str := C.GoString(C.olm_session_last_error(s.sess))
	if str == noErr {
		return nil
	}
	return errors.New("olm: session error: " + str)
}

// Clear clears the backing memory of Session.
func (s *Session) Clear() error {
	ret := C.olm_clear_session(s.sess)
	if ret == errValue {
		return s.LastError()
	}
	return nil
}

// ID returns the ID of the session. It will be the same on both ends.
func (s *Session) ID() (string, error) {
	length := C.olm_session_id_length(s.sess)
	id := C.malloc(length)
	defer C.free(id)

	ret := C.olm_session_id(s.sess, id, length)
	if ret == errValue {
		return "", s.LastError()
	}
	return C.GoString((*C.char)(id)), nil
}

// MatchesInbound returns if the provided OTK message is for this inbound session.
func (s *Session) MatchesInbound(otkMsg string) (bool, error) {
	cOTKMsg := unsafe.Pointer(C.CString(otkMsg))
	defer C.free(cOTKMsg)

	ret := C.olm_matches_inbound_session(s.sess, cOTKMsg, C.size_t(len(otkMsg)))
	if ret == errValue {
		return false, s.LastError()
	}
	return ret == 1, nil
}

// Pickle stores the session as a base64 string encrypted by key.
func (s *Session) Pickle(key string) (string, error) {
	length := C.olm_pickle_session_length(s.sess)
	pickle := C.malloc(length)
	cKey := unsafe.Pointer(C.CString(key))
	defer C.free(pickle)
	defer C.free(cKey)

	ret := C.olm_pickle_session(s.sess, cKey, C.size_t(len(key)), pickle, length)
	if ret == errValue {
		return "", s.LastError()
	}
	return C.GoStringN((*C.char)(pickle), C.int(ret)), nil
}

// Encrypt encrypts the provided plaintext and returns the message type with the encrypted message.
func (s *Session) Encrypt(plaintext string) (int, string, error) {
	randomLen := C.olm_encrypt_random_length(s.sess)
	randomBytes := make([]byte, randomLen)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return 0, "", fmt.Errorf("olm: error fetching random bytes: %w", err)
	}

	cBytes := C.CBytes(randomBytes)
	defer C.free(cBytes)

	msgType := C.olm_encrypt_message_type(s.sess)
	if msgType == errValue {
		return 0, "", s.LastError()
	}

	msgLen := C.olm_encrypt_message_length(s.sess, C.size_t(len(plaintext)))
	cPlaintext := unsafe.Pointer(C.CString(plaintext))
	cMsg := C.malloc(msgLen)

	ret := C.olm_encrypt(s.sess, cPlaintext, C.size_t(len(plaintext)), cBytes, randomLen, cMsg, msgLen)
	if ret == errValue {
		return 0, "", s.LastError()
	}
	return int(msgType), C.GoStringN((*C.char)(cMsg), C.int(ret)), nil
}

// Decrypt decrypts the provided encrypted text.
func (s *Session) Decrypt(msgType int, encrypted string) (string, error) {
	cEncrypted := unsafe.Pointer(C.CString(encrypted))
	defer C.free(cEncrypted)

	plaintextLen := C.olm_decrypt_max_plaintext_length(s.sess, C.size_t(msgType), cEncrypted, C.size_t(len(encrypted)))
	plaintext := C.malloc(plaintextLen)
	defer C.free(plaintext)

	ret := C.olm_decrypt(s.sess, C.size_t(msgType), cEncrypted, C.size_t(len(encrypted)), plaintext, plaintextLen)
	if ret == errValue {
		return "", s.LastError()
	}
	return C.GoStringN((*C.char)(plaintext), C.int(ret)), nil
}
