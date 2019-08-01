package log

import (
	"sync"

	"github.com/cultureamp/glamplify/config"
)

// LoggerFactory contains all the registered loggers
type LoggerFactory struct {
	loggers    map[string]ILogger
	nullLogger ILogger
}

// Factory to retrieve registered loggers
var (
	factory *LoggerFactory
	once    sync.Once
)

// Get the default logger, or if not set the nullLogger
func Get() ILogger {
	return GetFor("default")
}

// GetFor retrieves a registered logger by name
func GetFor(loggerName string) ILogger {
	once.Do(func() {
		factory = newFactory()
	})

	logger, ok := factory.loggers[loggerName]
	if !ok {
		return factory.nullLogger
	}

	return logger
}

func newFactory() *LoggerFactory {

	f := &LoggerFactory{}
	f.loggers = make(map[string]ILogger)

	// Create the default, NullLogger
	f.nullLogger = newNullLogger()

	settings := config.Load()
	for _, logConfig := range settings.App.Loggers {
		logger := newLogger(logConfig.Name, logConfig.Level)
		f.loggers[logConfig.Name] = logger
	}

	return f
}
