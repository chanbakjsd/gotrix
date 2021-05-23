package olm

// #include "olm/olm.h"
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

var util = newUtility()

type utility struct {
	mu   sync.Mutex
	util *C.OlmUtility
}

// newUtility creates a new Utility struct. It does not free as only one instance should be initialized at any point.
// Use the util global variable instead.
func newUtility() *utility {
	mem := C.malloc(utilitySize)
	return &utility{
		util: C.olm_utility(mem),
	}
}

// LastError returns the last error utility returned.
func (u *utility) LastError() error {
	str := C.GoString(C.olm_utility_last_error(u.util))
	if str == noErr {
		return nil
	}
	return errors.New("olm: utility error: " + str)
}

// Clear clears the backing memory of utility.
func (u *utility) Clear() error {
	ret := C.olm_clear_utility(u.util)
	if ret == errValue {
		return u.LastError()
	}
	return nil
}

// SHA256 calculates the base-64 encoded SHA256 hash.
func SHA256(input string) (string, error) {
	util.mu.Lock()
	defer util.mu.Unlock()

	length := C.olm_sha256_length(util.util)
	output := C.malloc(length)
	cInput := unsafe.Pointer(C.CString(input))

	ret := C.olm_sha256(util.util, cInput, C.size_t(len(input)), output, length)
	if ret == errValue {
		return "", util.LastError()
	}
	return C.GoStringN((*C.char)(output), C.int(length)), nil
}

// VerifyED25519 verifies that the message and signature is associated with key.
func VerifyED25519(key string, message string, signature string) error {
	util.mu.Lock()
	defer util.mu.Unlock()

	cKey := unsafe.Pointer(C.CString(key))
	cMessage := unsafe.Pointer(C.CString(message))
	cSignature := unsafe.Pointer(C.CString(signature))
	ret := C.olm_ed25519_verify(
		util.util, cKey, C.size_t(len(key)), cMessage, C.size_t(len(message)),
		cSignature, C.size_t(len(signature)),
	)
	if ret == errValue {
		return util.LastError()
	}
	return nil
}
