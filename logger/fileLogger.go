package logger

import (
	"io"
	"os"
	log "github.com/sirupsen/logrus"
)

type fileLogger struct {
	name   string
	logrus *log.Logger
}

// New AmplifyLogger
func newFileLogger(name string, formatter string, fullTimestamp bool, output string, level string) ILogger {

	return &fileLogger{
		name:   name,
		logrus: configureNewLogger(formatter, fullTimestamp, output, level),
	}
}

// Debugf todo
func (logger *fileLogger) Debugf(format string, args ...interface{}) {
	logger.logrus.Debugf(format, args...)
}

// Infof todo
func (logger *fileLogger) Infof(format string, args ...interface{}) {
	logger.logrus.Infof(format, args...)
}

// Warnf todo
func (logger *fileLogger) Warnf(format string, args ...interface{}) {
	logger.logrus.Warnf(format, args...)
}

// Errorf todo
func (logger *fileLogger) Errorf(format string, args ...interface{}) {
	logger.logrus.Errorf(format, args...)
}

// Fatalf todo
func (logger *fileLogger) Fatalf(format string, args ...interface{}) {
	logger.logrus.Fatalf(format, args...)
}

// Panicf todo
func (logger *fileLogger) Panicf(format string, args ...interface{}) {
	logger.logrus.Panicf(format, args...)
}

// Debug todo
func (logger *fileLogger) Debug(args ...interface{}) {
	logger.logrus.Debug(args...)
}

// Info todo
func (logger *fileLogger) Info(args ...interface{}) {
	logger.logrus.Info(args...)
}

// Warn todo
func (logger *fileLogger) Warn(args ...interface{}) {
	logger.logrus.Warn(args...)
}

// Error todo
func (logger *fileLogger) Error(args ...interface{}) {
	logger.logrus.Error(args...)
}

// Fatal todo
func (logger *fileLogger) Fatal(args ...interface{}) {
	logger.logrus.Fatal(args...)
}

// Panic todo
func (logger *fileLogger) Panic(args ...interface{}) {
	logger.logrus.Panic(args...)
}

// SetOutput todo
func (logger *fileLogger) SetOutput(output io.Writer) {
	logger.logrus.SetOutput(output)
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
