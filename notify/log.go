package notify

import (
	"fmt"

	"github.com/cultureamp/glamplify/log"
)

type notifyLogger struct {
	fieldLogger *log.FieldLogger
}

func newNotifyLogger() *notifyLogger {
	logger := log.New()

	return &notifyLogger{
		fieldLogger: logger,
	}
}

func (logger notifyLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.fieldLogger.Info(msg)
}
