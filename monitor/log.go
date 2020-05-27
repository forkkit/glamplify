package monitor

import (
	"context"
	"errors"

	"github.com/cultureamp/glamplify/log"
)

// Logger is the interface that is used for logging in the New Relic go-agent.  Assign the
// config.Logger types to the Logger you wish to use.  Loggers must be safe for
// use in multiple goroutines.
type monitorLogger struct {
	logger *log.Logger
}

func newMonitorLogger(ctx context.Context) *monitorLogger {
	writer := NewWriter()
	logger := log.NewFromCtxWithCustomerWriter(ctx, writer)

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