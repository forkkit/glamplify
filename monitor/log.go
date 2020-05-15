package monitor

import (
	"context"
	"errors"

	"github.com/cultureamp/glamplify/log"
)

// Logger is the interface that is used for logging in the New Relic go-agent.  Assign the
// config.Logger types to the Logger you wish to use.  Loggers must be safe for
// use in multiple goroutines.  Two Logger implementations are included:
// NewLogger, which logs at info level, and NewDebugLogger which logs at debug
// level.  logrus and logxi are supported by the integration packages
// https://godoc.org/github.com/newrelic/go-agent/_integrations/nrlogrus and
// https://godoc.org/github.com/newrelic/go-agent/_integrations/nrlogxi/v1.
type monitorLogger struct {
	logger *log.Logger
}

func newMonitorLogger(ctx context.Context) *monitorLogger {
	logger := log.NewFromCtx(ctx)

	return &monitorLogger{
		logger: logger,
	}
}

func (logger monitorLogger) Error(msg string, context map[string]interface{}) {
	err := errors.New(msg)
	logger.logger.Error("monitor_error", err, context)
}

func (logger monitorLogger) Warn(msg string, context map[string]interface{}) {
	logger.logger.Warn("monitor_warn", context, log.Fields{
		log.Message: msg,
	})
}

func (logger monitorLogger) Info(msg string, context map[string]interface{}) {
	logger.logger.Info("monitor_info", context, context, log.Fields{
		log.Message: msg,
	})
}

func (logger monitorLogger) Debug(msg string, context map[string]interface{}) {
	logger.logger.Debug("monitor_debug", context, context, log.Fields{
		log.Message: msg,
	})
}

func (logger monitorLogger) DebugEnabled() bool {
	return false
}