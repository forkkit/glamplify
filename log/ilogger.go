package log

// Level type
type Level uint32

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// ILogger interface that all loggers much support
type ILogger interface {

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
	DebugWithFields(fields Fields, args ...interface{})
	// Info logs info level logging
	InfoWithFields(fields Fields, args ...interface{})
	// Warn logs warn level logging
	WarnWithFields(fields Fields, args ...interface{})
	// Error logs error level logging
	ErrorWithFields(fields Fields, args ...interface{})
	// Fatal logs fatal level logging then the process will exit with status set to 1
	FatalWithFields(fields Fields, args ...interface{})
	// Panic logs panic level logging hen prints the stack trace and starts terminating the process unless recover is called.
	PanicWithFields(fields Fields, args ...interface{})

	// Debugf logs formatted debug level logging
	DebugfWithFields(fields Fields, format string, args ...interface{})
	// Infof logs formatted info level logging
	InfofWithFields(fields Fields, format string, args ...interface{})
	// Warnf logs formatted warn level logging
	WarnfWithFields(fields Fields, format string, args ...interface{})
	// Errorf logs formatted error level logging
	ErrorfWithFields(fields Fields, format string, args ...interface{})
	// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
	FatalfWithFields(fields Fields, format string, args ...interface{})
	// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
	PanicfWithFields(fields Fields, format string, args ...interface{})

	// SetOutput sets the output of the logger to the io.Writer
	AddHook(hook ILoggerHook)
}
