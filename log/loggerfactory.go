package log

import (
	"github.com/cultureamp/glamplify/config"
)

// Factory contains all the registered loggers
type LoggerFactory struct {
	loggers    map[string]ILogger
	nullLogger ILogger
}

// LoggerFactory to retrieve registered loggers
var Factory *LoggerFactory

// Get a registered logger by name
func (factory *LoggerFactory) Get(loggerName string) ILogger {

	logger, ok := factory.loggers[loggerName]
	if !ok {
		return factory.nullLogger
	}

	return logger
}

func init() {

	Factory = &LoggerFactory{}
	Factory.loggers = make(map[string]ILogger)

	// Create the default, NullLogger
	Factory.nullLogger = newNullLogger()

	for _, logConfig := range config.Settings.App.Loggers {
		logger := newLogger(logConfig.Name, logConfig.Level)
		Factory.loggers[logConfig.Name] = logger
	}
}
