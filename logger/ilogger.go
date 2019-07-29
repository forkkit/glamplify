package logger

import (
	"context"
	"runtime"
	"time"
)

// LogFields Fields type, used to pass to `WithFields`.
type LogFields map[string]interface{}

// LogLevel todo
type LogLevel uint32

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLogLevel LogLevel = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLogLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLogLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLogLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLogLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLogLevel
)

// AllLogLevels A constant exposing all logging levels
var AllLogLevels = []LogLevel{
	PanicLogLevel,
	FatalLogLevel,
	ErrorLogLevel,
	WarnLogLevel,
	InfoLogLevel,
	DebugLogLevel,
}

// An entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.
type LogEntry struct {

	// Contains all the fields set by the user.
	Fields LogFields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Trace, Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level LogLevel

	// Calling method, with package name
	Caller *runtime.Frame

	// Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// Contains the context set by the user. Useful for hook processing etc.
	Context context.Context
}

// ILoggerHook todo
type ILoggerHook interface {
	Levels() []LogLevel
	Fire(entry *LogEntry)
}

// ILogger interface that all loggers much support
type ILogger interface {

	// Debugf logs formatted debug level logging
	Debugf(format string, args ...interface{})
	// Infof logs formatted info level logging
	Infof(format string, args ...interface{})
	// Warnf logs formatted warn level logging
	Warnf(format string, args ...interface{})
	// Errorf logs formatted error level logging
	Errorf(format string, args ...interface{})
	// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
	Fatalf(format string, args ...interface{})
	// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
	Panicf(format string, args ...interface{})

	// Debug logs debug level logging
	Debug(args ...interface{})
	// Info logs info level logging
	Info(args ...interface{})
	// Warn logs warn level logging
	Warn(args ...interface{})
	// Error logs error level logging
	Error(args ...interface{})
	// Fatal logs fatal level logging then the process will exit with status set to 1
	Fatal(args ...interface{})
	// Panic logs panic level logging hen prints the stack trace and starts terminating the process unless recover is called.
	Panic(args ...interface{})

	// SetOutput sets the output of the logger to the io.Writer
	AddHook(hook ILoggerHook)
}
