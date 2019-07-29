package logger

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

type fileLogger struct {
	name   string
	logrus *log.Logger
	hooks  map[LogLevel][]ILoggerHook
	lock   sync.Mutex
}

// Debugf logs formatted debug level logging
func (logger *fileLogger) Debugf(format string, args ...interface{}) {
	logger.logrus.Debugf(format, args...)
}

// Infof logs formatted info level logging
func (logger *fileLogger) Infof(format string, args ...interface{}) {
	logger.logrus.Infof(format, args...)
}

// Warnf logs formatted warn level logging
func (logger *fileLogger) Warnf(format string, args ...interface{}) {
	logger.logrus.Warnf(format, args...)
}

// Errorf logs formatted error level logging
func (logger *fileLogger) Errorf(format string, args ...interface{}) {
	logger.logrus.Errorf(format, args...)
}

// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
func (logger *fileLogger) Fatalf(format string, args ...interface{}) {
	logger.logrus.Fatalf(format, args...)
}

// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
func (logger *fileLogger) Panicf(format string, args ...interface{}) {
	logger.logrus.Panicf(format, args...)
}

// Debug logs debug level logging
func (logger *fileLogger) Debug(args ...interface{}) {
	logger.logrus.Debug(args...)
}

// Info logs info level logging
func (logger *fileLogger) Info(args ...interface{}) {
	logger.logrus.Info(args...)
}

// Warn logs warn level logging
func (logger *fileLogger) Warn(args ...interface{}) {
	logger.logrus.Warn(args...)
}

// Error logs error level logging
func (logger *fileLogger) Error(args ...interface{}) {
	logger.logrus.Error(args...)
}

// Fatal logs fatal level logging then the process will exit with status set to 1
func (logger *fileLogger) Fatal(args ...interface{}) {
	logger.logrus.Fatal(args...)
}

// Panic logs panic level logging hen prints the stack trace and starts terminating the process unless recover is called.
func (logger *fileLogger) Panic(args ...interface{}) {
	logger.logrus.Panic(args...)
}

// SetOutput sets the output of the logger to the io.Writer
func (logger *fileLogger) AddHook(hook ILoggerHook) {
	logger.lock.Lock()
	defer logger.lock.Unlock()

	for _, level := range hook.Levels() {
		logger.hooks[level] = append(logger.hooks[level], hook)
	}
}

func (logger *fileLogger) Fire(entry *log.Entry) error {
	logEntry := convertEntryToLogEntry(entry)

	logger.lock.Lock()
	defer logger.lock.Unlock()

	for _, hook := range logger.hooks[logEntry.Level] {
		hook.Fire(logEntry)
	}

	return nil
}

func (logger *fileLogger) Levels() []log.Level {
	return log.AllLevels
}

func convertEntryToLogEntry(entry *log.Entry) *LogEntry {
	logEntry := &LogEntry{
		Time:    entry.Time,
		Level:   (LogLevel)(entry.Level),
		Caller:  entry.Caller,
		Message: entry.Message,
	}

	logEntry.Fields = convertDataToLogData(entry.Data)

	return logEntry
}

func convertDataToLogData(fields log.Fields) LogFields {

	logFields := make(LogFields)
	for k, v := range fields {
		logFields[k] = v
	}
	return logFields
}

func newFileLogger(name string, formatter string, fullTimestamp bool, output string, level string) ILogger {

	logger := &fileLogger{
		name:   name,
		logrus: configureNewLogger(formatter, fullTimestamp, output, level),
	}

	logger.hooks = make(map[LogLevel][]ILoggerHook)
	logger.logrus.AddHook(logger)
	return logger
}

func configureNewLogger(formatter string, fullTimestamp bool, output string, level string) *log.Logger {
	logger := log.New()

	if formatter == "json" {
		logger.SetFormatter(&log.JSONFormatter{})
	} else {
		logger.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: fullTimestamp,
		})
	}

	if output == "stdout" {
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(os.Stderr)
	}

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
