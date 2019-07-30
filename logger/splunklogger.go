package logger

import (
	"io/ioutil"
	"sync"

	splunk "github.com/Franco-Poveda/logrus-splunk-hook"
	"github.com/sirupsen/logrus"
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

func (logger *splunkLogger) Fire(entry *logrus.Entry) error {
	logEntry := convertEntryToLogEntry(entry)

	logger.lock.Lock()
	defer logger.lock.Unlock()

	for _, hook := range logger.hooks[logEntry.Level] {
		hook.Fire(logEntry)
	}

	return nil
}

func (logger *splunkLogger) Levels() []logrus.Level {
	return logrus.AllLevels
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
	client := &splunk.Client{} // <- !!!!
	h := splunk.NewHook(client, logrus.AllLevels)
	logger.logrus.AddHook(h)

	return logger
}

func configureNewSplunkLogger(formatter string, fullTimestamp bool, level string) *logrus.Logger {
	logger := logrus.New()

	if formatter == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: fullTimestamp,
		})
	}

	// This logger doesn't write anywhere, we just use its "hook" to send to splunk
	logger.SetOutput(ioutil.Discard)

	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.PanicLevel)
	}

	return logger
}
