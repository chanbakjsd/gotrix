package debug

import (
	"sync"

	"github.com/fatih/color"
)

var (
	traceColor   = color.New(color.FgCyan)
	debugColor   = color.New(color.FgHiCyan)
	infoColor    = color.New(color.FgHiBlue)
	warningColor = color.New(color.FgHiYellow)
	errorColor   = color.New(color.FgHiRed)

	// There's potentially multiple input source from different goroutines.
	// Namely, each event handler is its own goroutine.
	logMutex sync.Mutex
)

type defaultLogger struct{}

func (defaultLogger) Trace(a interface{}) {
	if TraceEnabled {
		logMutex.Lock()
		defer logMutex.Unlock()
		_, _ = traceColor.Println(a)
	}
}
func (defaultLogger) Debug(a interface{}) {
	if DebugEnabled {
		logMutex.Lock()
		defer logMutex.Unlock()
		_, _ = debugColor.Println(a)
	}
}
func (defaultLogger) Info(a interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	_, _ = infoColor.Println(a)
}
func (defaultLogger) Warn(a interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	_, _ = warningColor.Println(a)
}
func (defaultLogger) Error(a interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	_, _ = errorColor.Println(a)
}
