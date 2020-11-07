package debug

import (
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	logrusUnderlyingLogger
}

type logrusUnderlyingLogger interface {
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})

	WithFields(logrus.Fields) *logrus.Entry
}

// NewLogrus creates a new Logger that uses Logrus as its underlying implementation.
func NewLogrus(level logrus.Level) LoggerType {
	logger := logrus.New()
	logger.Level = level
	return logrusLogger{logger}
}

func (l logrusLogger) Fields(a map[string]interface{}) LoggerType {
	return logrusLogger{
		l.WithFields(logrus.Fields(a)),
	}
}
