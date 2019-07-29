package logger

import "io"

type nullLogger struct {
	name string
}

// New AmplifyLogger
func newNullLogger() ILogger {

	return &nullLogger{
		name: "<null>",
	}
}

// Debugf todo
func (logger *nullLogger) Debugf(format string, args ...interface{}) {
}

// Infof todo
func (logger *nullLogger) Infof(format string, args ...interface{}) {
}

// Warnf todo
func (logger *nullLogger) Warnf(format string, args ...interface{}) {
}

// Errorf todo
func (logger *nullLogger) Errorf(format string, args ...interface{}) {
}

// Fatalf todo
func (logger *nullLogger) Fatalf(format string, args ...interface{}) {
}

// Panicf todo
func (logger *nullLogger) Panicf(format string, args ...interface{}) {
}

// Debug todo
func (logger *nullLogger) Debug(args ...interface{}) {
}

// Info todo
func (logger *nullLogger) Info(args ...interface{}) {
}

// Warn todo
func (logger *nullLogger) Warn(args ...interface{}) {
}

// Error todo
func (logger *nullLogger) Error(args ...interface{}) {
}

// Fatal todo
func (logger *nullLogger) Fatal(args ...interface{}) {
}

// Panic todo
func (logger *nullLogger) Panic(args ...interface{}) {
}

// SetOutput todo
func (logger *nullLogger) SetOutput(output io.Writer) {
}
