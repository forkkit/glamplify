package log

import (
	"github.com/cultureamp/glamplify/config"
)

// Factory contains all the registered loggers
type Factory struct {
	loggers    map[string]ILogger
	nullLogger ILogger
}

// LoggerFactory to retrieve registered loggers
var LoggerFactory *Factory

// Get a registered logger by name
func (factory *Factory) Get(loggerName string) ILogger {

	logger, ok := factory.loggers[loggerName]
	if !ok {
		return factory.nullLogger
	}

	return logger
}

func init() {

	LoggerFactory = &Factory{}
	LoggerFactory.loggers = make(map[string]ILogger)

	// Create the default, NullLogger
	LoggerFactory.nullLogger = newNullLogger()

	for _, logConfig := range config.Settings.App.Loggers {
		logger := newLogger(logConfig.Name, logConfig.Level)
		LoggerFactory.loggers[logConfig.Name] = logger
	}
}
