package logger

import (
	"io/ioutil"
	"sync"

	"github.com/multiplay/go-slack/chat"
	"github.com/multiplay/go-slack/lrhook"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type slackLogger struct {
	name    string
	url     string
	channel string
	emoji   string

	logrus *log.Logger
	hooks  map[LogLevel][]ILoggerHook
	lock   sync.Mutex
}

// Debugf logs formatted debug level logging
func (logger *slackLogger) Debugf(format string, args ...interface{}) {
	logger.logrus.Debugf(format, args...)
}

// Infof logs formatted info level logging
func (logger *slackLogger) Infof(format string, args ...interface{}) {
	logger.logrus.Infof(format, args...)
}

// Warnf logs formatted warn level logging
func (logger *slackLogger) Warnf(format string, args ...interface{}) {
	logger.logrus.Warnf(format, args...)
}

// Errorf logs formatted error level logging
func (logger *slackLogger) Errorf(format string, args ...interface{}) {
	logger.logrus.Errorf(format, args...)
}

// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
func (logger *slackLogger) Fatalf(format string, args ...interface{}) {
	logger.logrus.Fatalf(format, args...)
}

// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
func (logger *slackLogger) Panicf(format string, args ...interface{}) {
	logger.logrus.Panicf(format, args...)
}

// Debug logs debug level logging
func (logger *slackLogger) Debug(args ...interface{}) {
	logger.logrus.Debug(args...)
}

// Info logs info level logging
func (logger *slackLogger) Info(args ...interface{}) {
	logger.logrus.Info(args...)
}

// Warn logs warn level logging
func (logger *slackLogger) Warn(args ...interface{}) {
	logger.logrus.Warn(args...)
}

// Error logs error level logging
func (logger *slackLogger) Error(args ...interface{}) {
	logger.logrus.Error(args...)
}

// Fatal logs fatal level logging then the process will exit with status set to 1
func (logger *slackLogger) Fatal(args ...interface{}) {
	logger.logrus.Fatal(args...)
}

// Panic logs panic level logging hen prints the stack trace and starts terminating the process unless recover is called.
func (logger *slackLogger) Panic(args ...interface{}) {
	logger.logrus.Panic(args...)
}

// SetOutput sets the output of the logger to the io.Writer
func (logger *slackLogger) AddHook(hook ILoggerHook) {
	logger.lock.Lock()
	defer logger.lock.Unlock()

	for _, level := range hook.Levels() {
		logger.hooks[level] = append(logger.hooks[level], hook)
	}
}

func (logger *slackLogger) Fire(entry *log.Entry) error {
	logEntry := convertEntryToLogEntry(entry)

	logger.lock.Lock()
	defer logger.lock.Unlock()

	for _, hook := range logger.hooks[logEntry.Level] {
		hook.Fire(logEntry)
	}

	return nil
}

func (logger *slackLogger) Levels() []log.Level {
	return log.AllLevels
}

func newSlackLogger(name string, formatter string, fullTimestamp bool, url string, channel string, emoji string, level string) ILogger {

	logger := &slackLogger{
		name:    name,
		url:     url,
		channel: channel,
		emoji:   emoji,
		logrus:  configureNewSlackLogger(formatter, fullTimestamp, level),
	}

	cfg := lrhook.Config{
		MinLevel: logrus.ErrorLevel,
		Message: chat.Message{
			Channel:   channel,
			IconEmoji: emoji,
		},
	}
	logger.hooks = make(map[LogLevel][]ILoggerHook)
	logger.logrus.AddHook(logger)

	h := lrhook.New(cfg, url)
	logger.logrus.AddHook(h)

	return logger
}

func configureNewSlackLogger(formatter string, fullTimestamp bool, level string) *log.Logger {
	logger := log.New()

	if formatter == "json" {
		logger.SetFormatter(&log.JSONFormatter{})
	} else {
		logger.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: fullTimestamp,
		})
	}

	// This logger doesn't write anywhere, we just use its "hook" to send to slack
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
