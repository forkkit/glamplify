package log

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/cultureamp/glamplify/config"
)

// LoggerFactory contains all the registered loggers
type LoggerFactory struct {
	loggers    map[string]*Logger
	nullLogger *Logger
}

// Factory to retrieve registered loggers
var (
	factory *LoggerFactory
	once    sync.Once
)

// Get the default logger, or if not set the nullLogger
func Get() *Logger {
	return GetFor("default")
}

// GetFor retrieves a registered logger by name
func GetFor(loggerName string) *Logger {
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
	f.loggers = make(map[string]*Logger)

	// Create the default, NullLogger
	f.nullLogger = newLogger("null", ioutil.Discard, "")

	settings := config.Load()
	for _, logConfig := range settings.App.Loggers {
		logger := newLogger(logConfig.Name, os.Stdout, logConfig.Level)
		f.loggers[logConfig.Name] = logger
	}

	return f
}
