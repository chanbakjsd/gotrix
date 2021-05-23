package olm

// #cgo LDFLAGS: -lolm
// #include <stdlib.h>
// #include "olm/olm.h"
import "C"

var (
	errValue = C.olm_error()
	noErr    = "SUCCESS"
)

// Sizes of the structs provided by olm.
var (
	accountSize = C.olm_account_size()
	sessionSize = C.olm_session_size()
	utilitySize = C.olm_utility_size()
)

// LibraryVersion returns the version of the cgo library embedded in the form of a major, minor, patch triplet.
func LibraryVersion() (uint8, uint8, uint8) {
	var major, minor, patch C.uchar
	C.olm_get_library_version(&major, &minor, &patch)
	return uint8(major), uint8(minor), uint8(patch)
}
