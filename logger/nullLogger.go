package logger

type nullLogger struct {
	name string
}

// Debugf logs formatted debug level logging
func (logger *nullLogger) Debugf(format string, args ...interface{}) {
}

// Infof logs formatted info level logging
func (logger *nullLogger) Infof(format string, args ...interface{}) {
}

// Warnf logs formatted warn level logging
func (logger *nullLogger) Warnf(format string, args ...interface{}) {
}

// Errorf logs formatted error level logging
func (logger *nullLogger) Errorf(format string, args ...interface{}) {
}

// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
func (logger *nullLogger) Fatalf(format string, args ...interface{}) {
}

// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
func (logger *nullLogger) Panicf(format string, args ...interface{}) {
}

// Debug logs debug level logging
func (logger *nullLogger) Debug(args ...interface{}) {
}

// Info logs info level logging
func (logger *nullLogger) Info(args ...interface{}) {
}

// Warn logs warn level logging
func (logger *nullLogger) Warn(args ...interface{}) {
}

// Error logs error level logging
func (logger *nullLogger) Error(args ...interface{}) {
}

// Fatal logs fatal level logging then the process will exit with status set to 1
func (logger *nullLogger) Fatal(args ...interface{}) {
}

// Panic logs panic level logging hen prints the stack trace and starts terminating the process unless recover is called.
func (logger *nullLogger) Panic(args ...interface{}) {
}

// SetOutput sets the output of the logger to the io.Writer
func (logger *nullLogger) AddHook(hook ILoggerHook) {
	// TODO - nullLogger does not honor hooks!
}

func newNullLogger() ILogger {

	return &nullLogger{
		name: "<null>",
	}
}
