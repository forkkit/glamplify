package log

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type logger struct {
	name  string
	level string

	logrus *logrus.Logger

	hooks []ILoggerHook
	hlock sync.Mutex
}

// Debug logs debug level logging
func (log *logger) Debug(args ...interface{}) {
	log.logrus.Debug(args...)
}

// Info logs info level logging
func (log *logger) Info(args ...interface{}) {
	log.logrus.Info(args...)
}

// Warn logs warn level logging
func (log *logger) Warn(args ...interface{}) {
	log.logrus.Warn(args...)
}

// Error logs error level logging
func (log *logger) Error(args ...interface{}) {
	log.logrus.Error(args...)
}

// Fatal logs fatal level logging then the process will exit with status set to 1
func (log *logger) Fatal(args ...interface{}) {
	log.logrus.Fatal(args...)
}

// Panic logs panic level logging hen prints the stack trace and call panic checking whether recover is called.
func (log *logger) Panic(args ...interface{}) {
	log.logrus.Panic(args...)
}

// Debugf logs formatted debug level logging
func (log *logger) Debugf(format string, args ...interface{}) {
	log.logrus.Debugf(format, args...)
}

// Infof logs formatted info level logging
func (log *logger) Infof(format string, args ...interface{}) {
	log.logrus.Infof(format, args...)
}

// Warnf logs formatted warn level logging
func (log *logger) Warnf(format string, args ...interface{}) {
	log.logrus.Warnf(format, args...)
}

// Errorf logs formatted error level logging
func (log *logger) Errorf(format string, args ...interface{}) {
	log.logrus.Errorf(format, args...)
}

// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
func (log *logger) Fatalf(format string, args ...interface{}) {
	log.logrus.Fatalf(format, args...)
}

// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
func (log *logger) Panicf(format string, args ...interface{}) {
	log.logrus.Panicf(format, args...)
}

// Debug logs debug level logging
func (log *logger) DebugWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Debug(args...)
}

// Info logs info level logging
func (log *logger) InfoWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Info(args...)
}

// Warn logs warn level logging
func (log *logger) WarnWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Warn(args...)
}

// Error logs error level logging
func (log *logger) ErrorWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Error(args...)
}

// Fatal logs fatal level logging then the process will exit with status set to 1
func (log *logger) FatalWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Fatal(args...)
}

// Panic logs panic level logging hen prints the stack trace and call panic checking whether recover is called.
func (log *logger) PanicWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Panic(args...)
}

// Debugf logs formatted debug level logging
func (log *logger) DebugfWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Debugf(format, args...)
}

// Infof logs formatted info level logging
func (log *logger) InfofWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Infof(format, args...)
}

// Warnf logs formatted warn level logging
func (log *logger) WarnfWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Warnf(format, args...)
}

// Errorf logs formatted error level logging
func (log *logger) ErrorfWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Errorf(format, args...)
}

// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
func (log *logger) FatalfWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Fatalf(format, args...)
}

// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
func (log *logger) PanicfWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Panicf(format, args...)
}

func (log *logger) AddHook(hook ILoggerHook) {
	log.hlock.Lock()
	defer log.hlock.Unlock()

	log.hooks = append(log.hooks, hook)
}

func (log *logger) Fire(entry *logrus.Entry) error {
	logEntry := convertEntry(entry)

	log.hlock.Lock()
	defer log.hlock.Unlock()

	for _, hook := range log.hooks {
		if log.logrus.IsLevelEnabled(entry.Level) {
			hook.Fire(logEntry)
		}
	}

	return nil
}

func (log *logger) Levels() []logrus.Level {
	return logrus.AllLevels
}

func newLogger(name string, level string) ILogger {

	logger := &logger{
		name:   name,
		logrus: newLogrusLogger(level),
	}

	logger.logrus.AddHook(logger)
	return logger
}

func newLogrusLogger(level string) *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		QuoteEmptyFields: true,
		FullTimestamp:    true,
	})

	// 12-Factor Apps always log to stdout
	logger.SetOutput(os.Stdout)

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

// Common logging routines
func convertEntry(entry *logrus.Entry) *Entry {
	e := &Entry{
		Time:    entry.Time,
		Level:   (Level)(entry.Level),
		Caller:  entry.Caller,
		Message: entry.Message,
	}

	e.Fields = convertFieldsFromLogrus(entry.Data)
	return e
}

func convertFieldsFromLogrus(fields logrus.Fields) Fields {

	Fields := make(Fields)
	for k, v := range fields {
		Fields[k] = v
	}
	return Fields
}

func convertFieldsToLogrus(Fields Fields) logrus.Fields {
	fields := make(logrus.Fields)

	for k, v := range Fields {
		fields[k] = v
	}

	return fields
}
