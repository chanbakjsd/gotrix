package debug

import (
	"log"
	"os"
)

var (
	// DebugEnabled is true if debug messages should be sent.
	DebugEnabled bool
	// TraceEnabled is true if messages should be traced.
	TraceEnabled bool
)

func init() {
	_, DebugEnabled = os.LookupEnv("GOMATRIX_DEBUG")
	_, TraceEnabled = os.LookupEnv("GOMATRIX_TRACE")

	DebugEnabled = DebugEnabled || TraceEnabled
}

// Debug prints the provided message out if Debug is true.
func Debug(a ...interface{}) {
	if DebugEnabled {
		log.Println(a...)
	}
}

// Trace prints the provided message out if Trace is true.
func Trace(a ...interface{}) {
	if TraceEnabled {
		log.Println(a...)
	}
}
