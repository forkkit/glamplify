package event

import (
	"errors"

	"github.com/cultureamp/glamplify/log"
)

// Logger is the interface that is used for logging in the go-agent.  Assign the
// Config.Logger field to the Logger you wish to use.  Loggers must be safe for
// use in multiple goroutines.  Two Logger implementations are included:
// NewLogger, which logs at info level, and NewDebugLogger which logs at debug
// level.  logrus and logxi are supported by the integration packages
// https://godoc.org/github.com/newrelic/go-agent/_integrations/nrlogrus and
// https://godoc.org/github.com/newrelic/go-agent/_integrations/nrlogxi/v1.
type eventLogger struct {
	fieldLogger *log.FieldLogger
}

func newEventLogger() *eventLogger {
	logger := log.New()

	return &eventLogger{
		fieldLogger: logger,
	}
}

func (logger eventLogger) Error(msg string, context map[string]interface{}) {
	err := errors.New(msg)
	logger.fieldLogger.Error(err, context)
}

func (logger eventLogger) Warn(msg string, context map[string]interface{}) {
	logger.fieldLogger.Print(msg, context)
}

func (logger eventLogger) Info(msg string, context map[string]interface{}) {
	logger.fieldLogger.Print(msg, context)
}

func (logger eventLogger) Debug(msg string, context map[string]interface{}) {
	logger.fieldLogger.Debug(msg, context)
}

func (logger eventLogger) DebugEnabled() bool {
	return false
}

func (logger eventLogger) merge(fields log.Fields, entries ...Entries) log.Fields {
	merged := log.Fields{}

	for k, v := range fields {
		merged[k] = v
	}

	for _, f := range entries {
		for k, v := range f {
			merged[k] = v
		}
	}

	return merged
}
