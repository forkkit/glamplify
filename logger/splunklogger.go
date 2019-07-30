package logger

import (
	"io/ioutil"
	"sync"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type splunkLogger struct {
	name       string
	url        string
	token      string
	source     string
	sourceType string
	index      string

	logrus *logrus.Logger
	hooks  map[LogLevel][]ILoggerHook
	lock   sync.Mutex
}

// Debugf logs formatted debug level logging
func (logger *splunkLogger) Debugf(format string, args ...interface{}) {
	logger.logrus.Debugf(format, args...)
}

// Infof logs formatted info level logging
func (logger *splunkLogger) Infof(format string, args ...interface{}) {
	logger.logrus.Infof(format, args...)
}

// Warnf logs formatted warn level logging
func (logger *splunkLogger) Warnf(format string, args ...interface{}) {
	logger.logrus.Warnf(format, args...)
}

// Errorf logs formatted error level logging
func (logger *splunkLogger) Errorf(format string, args ...interface{}) {
	logger.logrus.Errorf(format, args...)
}

// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
func (logger *splunkLogger) Fatalf(format string, args ...interface{}) {
	logger.logrus.Fatalf(format, args...)
}

// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
func (logger *splunkLogger) Panicf(format string, args ...interface{}) {
	logger.logrus.Panicf(format, args...)
}

// Debug logs debug level logging
func (logger *splunkLogger) Debug(args ...interface{}) {
	logger.logrus.Debug(args...)
}

// Info logs info level logging
func (logger *splunkLogger) Info(args ...interface{}) {
	logger.logrus.Info(args...)
}

// Warn logs warn level logging
func (logger *splunkLogger) Warn(args ...interface{}) {
	logger.logrus.Warn(args...)
}

// Error logs error level logging
func (logger *splunkLogger) Error(args ...interface{}) {
	logger.logrus.Error(args...)
}

// Fatal logs fatal level logging then the process will exit with status set to 1
func (logger *splunkLogger) Fatal(args ...interface{}) {
	logger.logrus.Fatal(args...)
}

// Panic logs panic level logging hen prints the stack trace and starts terminating the process unless recover is called.
func (logger *splunkLogger) Panic(args ...interface{}) {
	logger.logrus.Panic(args...)
}

// SetOutput sets the output of the logger to the io.Writer
func (logger *splunkLogger) AddHook(hook ILoggerHook) {
	logger.lock.Lock()
	defer logger.lock.Unlock()

	for _, level := range hook.Levels() {
		logger.hooks[level] = append(logger.hooks[level], hook)
	}
}

func (logger *splunkLogger) Fire(entry *log.Entry) error {
	logEntry := convertEntryToLogEntry(entry)

	logger.lock.Lock()
	defer logger.lock.Unlock()

	for _, hook := range logger.hooks[logEntry.Level] {
		hook.Fire(logEntry)
	}

	return nil
}

func (logger *splunkLogger) Levels() []log.Level {
	return log.AllLevels
}

func newSplunkLogger(name string, formatter string, fullTimestamp bool, url string, token string, source string, sourceType string, index string, level string) ILogger {

	logger := &splunkLogger{
		name:       name,
		url:        url,
		token:      token,
		source:     source,
		sourceType: sourceType,
		index:      index,
		logrus:     configureNewSplunkLogger(formatter, fullTimestamp, level),
	}

	logger.hooks = make(map[LogLevel][]ILoggerHook)
	logger.logrus.AddHook(logger)

	// TODO
	//	h := lrhook.New(cfg, url)
	//	logger.logrus.AddHook(h)

	return logger
}

func configureNewSplunkLogger(formatter string, fullTimestamp bool, level string) *log.Logger {
	logger := log.New()

	if formatter == "json" {
		logger.SetFormatter(&log.JSONFormatter{})
	} else {
		logger.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: fullTimestamp,
		})
	}

	// This logger doesn't write anywhere, we just use its "hook" to send to splunk
	logger.SetOutput(ioutil.Discard)

	switch level {
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	case "fatal":
		logger.SetLevel(log.FatalLevel)
	default:
		logger.SetLevel(log.PanicLevel)
	}

	return logger
}
