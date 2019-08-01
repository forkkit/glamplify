package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	Name   string
	output io.Writer
	Level  string

	logrus *logrus.Logger
	hooks  *Hook
}

// Debug logs debug level logging
func (log *Logger) Debug(args ...interface{}) {
	log.logrus.Debug(args...)
}

// Info logs info level logging
func (log *Logger) Info(args ...interface{}) {
	log.logrus.Info(args...)
}

// Warn logs warn level logging
func (log *Logger) Warn(args ...interface{}) {
	log.logrus.Warn(args...)
}

// Error logs error level logging
func (log *Logger) Error(args ...interface{}) {
	log.logrus.Error(args...)
}

// Fatal logs fatal level logging then the process will exit with status set to 1
func (log *Logger) Fatal(args ...interface{}) {
	log.logrus.Fatal(args...)
}

// Panic logs panic level logging hen prints the stack trace and call panic checking whether recover is called.
func (log *Logger) Panic(args ...interface{}) {
	log.logrus.Panic(args...)
}

// Debugf logs formatted debug level logging
func (log *Logger) Debugf(format string, args ...interface{}) {
	log.logrus.Debugf(format, args...)
}

// Infof logs formatted info level logging
func (log *Logger) Infof(format string, args ...interface{}) {
	log.logrus.Infof(format, args...)
}

// Warnf logs formatted warn level logging
func (log *Logger) Warnf(format string, args ...interface{}) {
	log.logrus.Warnf(format, args...)
}

// Errorf logs formatted error level logging
func (log *Logger) Errorf(format string, args ...interface{}) {
	log.logrus.Errorf(format, args...)
}

// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
func (log *Logger) Fatalf(format string, args ...interface{}) {
	log.logrus.Fatalf(format, args...)
}

// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
func (log *Logger) Panicf(format string, args ...interface{}) {
	log.logrus.Panicf(format, args...)
}

// DebugWithFields logs debug level logging
func (log *Logger) DebugWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Debug(args...)
}

// InfoWithFields logs info level logging
func (log *Logger) InfoWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Info(args...)
}

// WarnWithFields logs warn level logging
func (log *Logger) WarnWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Warn(args...)
}

// ErrorWithFields logs error level logging
func (log *Logger) ErrorWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Error(args...)
}

// FatalWithFields logs fatal level logging then the process will exit with status set to 1
func (log *Logger) FatalWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Fatal(args...)
}

// PanicWithFields logs panic level logging hen prints the stack trace and call panic checking whether recover is called.
func (log *Logger) PanicWithFields(fields Fields, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Panic(args...)
}

// DebugfWithFields logs formatted debug level logging
func (log *Logger) DebugfWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Debugf(format, args...)
}

// InfofWithFields logs formatted info level logging
func (log *Logger) InfofWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Infof(format, args...)
}

// WarnfWithFields logs formatted warn level logging
func (log *Logger) WarnfWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Warnf(format, args...)
}

// ErrorfWithFields logs formatted error level logging
func (log *Logger) ErrorfWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Errorf(format, args...)
}

// FatalfWithFields logs formatted fatal level logging then the process will exit with status set to 1
func (log *Logger) FatalfWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Fatalf(format, args...)
}

// PanicfWithFields logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
func (log *Logger) PanicfWithFields(fields Fields, format string, args ...interface{}) {
	f := convertFieldsToLogrus(fields)
	log.logrus.WithFields(f).Panicf(format, args...)
}

// AddHook todo...
func (log *Logger) AddHook(hook ILoggerHook) {
	log.hooks.AddHook(hook)
}

func newLogger(name string, output io.Writer, level string) *Logger {

	logger := &Logger{
		Name:   name,
		output: output,
		Level:  level,
		logrus: newLogrusLogger(output, level),
	}

	logger.hooks = newHook(logger)
	return logger
}

func newLogrusLogger(output io.Writer, level string) *logrus.Logger {
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

func convertFieldsToLogrus(Fields Fields) logrus.Fields {
	fields := make(logrus.Fields)

	for k, v := range Fields {
		fields[k] = v
	}

	return fields
}
