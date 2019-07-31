package log

type nullLogger struct {
	name string
}

// Debugf logs formatted debug level logging
func (log *nullLogger) Debugf(format string, args ...interface{}) {
}

// Infof logs formatted info level logging
func (log *nullLogger) Infof(format string, args ...interface{}) {
}

// Warnf logs formatted warn level logging
func (log *nullLogger) Warnf(format string, args ...interface{}) {
}

// Errorf logs formatted error level logging
func (log *nullLogger) Errorf(format string, args ...interface{}) {
}

// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
func (log *nullLogger) Fatalf(format string, args ...interface{}) {
}

// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
func (log *nullLogger) Panicf(format string, args ...interface{}) {
}

// Debug logs debug level logging
func (log *nullLogger) Debug(args ...interface{}) {
}

// Info logs info level logging
func (log *nullLogger) Info(args ...interface{}) {
}

// Warn logs warn level logging
func (log *nullLogger) Warn(args ...interface{}) {
}

// Error logs error level logging
func (log *nullLogger) Error(args ...interface{}) {
}

// Fatal logs fatal level logging then the process will exit with status set to 1
func (log *nullLogger) Fatal(args ...interface{}) {
}

// Panic logs panic level logging hen prints the stack trace and starts terminating the process unless recover is called.
func (log *nullLogger) Panic(args ...interface{}) {
}

// Debug logs debug level logging
func (log *nullLogger) DebugWithFields(fields Fields, args ...interface{}) {
}

// Info logs info level logging
func (log *nullLogger) InfoWithFields(fields Fields, args ...interface{}) {
}

// Warn logs warn level logging
func (log *nullLogger) WarnWithFields(fields Fields, args ...interface{}) {
}

// Error logs error level logging
func (log *nullLogger) ErrorWithFields(fields Fields, args ...interface{}) {
}

// Fatal logs fatal level logging then the process will exit with status set to 1
func (log *nullLogger) FatalWithFields(fields Fields, args ...interface{}) {
}

// Panic logs panic level logging hen prints the stack trace and call panic checking whether recover is called.
func (log *nullLogger) PanicWithFields(fields Fields, args ...interface{}) {
}

// Debugf logs formatted debug level logging
func (log *nullLogger) DebugfWithFields(fields Fields, format string, args ...interface{}) {
}

// Infof logs formatted info level logging
func (log *nullLogger) InfofWithFields(fields Fields, format string, args ...interface{}) {
}

// Warnf logs formatted warn level logging
func (log *nullLogger) WarnfWithFields(fields Fields, format string, args ...interface{}) {
}

// Errorf logs formatted error level logging
func (log *nullLogger) ErrorfWithFields(fields Fields, format string, args ...interface{}) {
}

// Fatalf logs formatted fatal level logging then the process will exit with status set to 1
func (log *nullLogger) FatalfWithFields(fields Fields, format string, args ...interface{}) {
}

// Panicf logs formatted panic level logging then prints the stack trace and starts terminating the process unless recover is called
func (log *nullLogger) PanicfWithFields(fields Fields, format string, args ...interface{}) {
}

// SetOutput sets the output of the logger to the io.Writer
func (log *nullLogger) AddHook(hook ILoggerHook) {
	// TODO - nullLogger does not honor hooks!
}

func newNullLogger() ILogger {

	return &nullLogger{
		name: "<null>",
	}
}
