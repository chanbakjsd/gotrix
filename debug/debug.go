package debug

import (
	"os"
)

// LoggerType is the type required by the Logger variable.
type LoggerType interface {
	Trace(interface{})
	Debug(interface{})
	Info(interface{})
	Warn(interface{})
	Error(interface{})
}

var (
	// DebugEnabled is true if debug messages should be sent.
	DebugEnabled bool
	// TraceEnabled is true if messages should be traced.
	TraceEnabled bool
)

// Logger is the default logger called.
// You can set it to redirect logs to other places.
var Logger LoggerType = defaultLogger{}

func init() {
	_, DebugEnabled = os.LookupEnv("GOMATRIX_DEBUG")
	_, TraceEnabled = os.LookupEnv("GOMATRIX_TRACE")

	DebugEnabled = DebugEnabled || TraceEnabled
}

// Trace calls Trace on the default Logger.
func Trace(a interface{}) {
	Logger.Trace(a)
}

// Debug calls Debug on the default Logger.
func Debug(a interface{}) {
	Logger.Debug(a)
}

// Info calls Info on the default Logger.
func Info(a interface{}) {
	Logger.Info(a)
}

// Warn calls Warn on the default Logger.
func Warn(a interface{}) {
	Logger.Warn(a)
}

// Error calls Error on the default Logger.
func Error(a interface{}) {
	Logger.Error(a)
}
