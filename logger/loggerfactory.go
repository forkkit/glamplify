package logger

import (
	. "github.com/cultureamp/gamplify/config"
)

// LogFactory contains all the registered loggers
type LogFactory struct {
	loggers    map[string]ILogger
	nullLogger ILogger
}

// LoggerFactory to retrieve registered loggers
var LoggerFactory *LogFactory

// Get a registered logger by name
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
	/*
		for _, _ = range Config.App.Loggers.Rules {

		}
	*/

	for _, _ = range Config.App.Loggers.Targets.Stream {

	}

	for _, _ = range Config.App.Loggers.Targets.Slack {

	}

	for _, _ = range Config.App.Loggers.Targets.Splunk {

	}

	/*
		for _, target := range Config.App.Loggers.Targets.StreamTargets {

			logger := newStreamLogger(
				target.Name,
				target.Formatter,
				target.FullTimestamp,
				target.Output,
				target.Level,
			)

			LoggerFactory.loggers[target.Name] = logger
		}
	*/
}
