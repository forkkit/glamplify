package logger

import (
	. "github.com/cultureamp/gamplify/config"
)

// LogFactory todo
type LogFactory struct {
	loggers    map[string]ILogger
	nullLogger ILogger
}

// LoggerFactory todo...
var LoggerFactory *LogFactory

// Get todo...
func (factory *LogFactory) Get(loggerName string) ILogger {

	logger, ok := factory.loggers[loggerName]
	if !ok {
		return factory.nullLogger
	}

	return logger
}

func init() {

	LoggerFactory = &LogFactory{}
	LoggerFactory.loggers = make(map[string]ILogger)

	// Create the default, NullLogger
	LoggerFactory.nullLogger = newNullLogger()

	// Loop through all the logs in the config and create specific loggers and add them to the map
	for _, target := range Config.App.Loggers {

		logger := newFileLogger(
			target.Name,
			target.Formatter,
			target.FullTimestamp,
			target.Output,
			target.Level,
		)

		LoggerFactory.loggers[target.Name] = logger
	}
}
