package log

import (
	"context"
	"net/http"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/helper"
)

// Logger
type Logger struct {
	rsFields  gcontext.RequestScopedFields
	fields    Fields
	sysValues *SystemValues
	writer    Writer
}

var (
	sevLevel = newSystemLogLevel()
	internalWriter = NewWriter(func(conf *WriterConfig) {})
	defaultLogger  = NewWitCustomWriter(gcontext.RequestScopedFields{}, internalWriter)
)

// New creates a *Logger with optional fields. Useful for when you want to add a field to all subsequent logging calls eg. request_id, etc.
func New(rsFields gcontext.RequestScopedFields, fields ...Fields) *Logger {
	return newLogger(rsFields, internalWriter, fields...)
}

// Useful for CLI applications that want to write to stderr or file etc.
func NewWitCustomWriter(rsFields gcontext.RequestScopedFields, writer Writer, fields ...Fields) *Logger {
	return newLogger(rsFields, writer, fields...)
}

// NewFromCtx creates a new logger from a context, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromCtx(ctx context.Context, fields ...Fields) *Logger {
	rsFields, _ := gcontext.GetRequestScopedFields(ctx)
	return New(rsFields, fields...)
}

// NewFromCtxWithCustomerWriter creates a new logger from a context, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromCtxWithCustomerWriter(ctx context.Context, writer Writer, fields ...Fields) *Logger {
	rsFields, _ := gcontext.GetRequestScopedFields(ctx)
	return NewWitCustomWriter(rsFields, writer, fields...)
}

// NewFromRequest creates a new logger from a http.Request, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromRequest(r *http.Request, fields ...Fields) *Logger {
	return NewFromCtx(r.Context(), fields...)
}

// NewFromRequest creates a new logger from a http.Request, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromRequestWithCustomWriter(r *http.Request, writer Writer, fields ...Fields) *Logger {
	return NewFromCtxWithCustomerWriter(r.Context(), writer, fields...)
}

func newLogger(rsFields gcontext.RequestScopedFields, writer Writer, fields ...Fields) *Logger {

	df := newSystemValues()

	merged := Fields{}
	merged = merged.Merge(fields...)
	logger := &Logger{
		rsFields: rsFields,
		writer:   writer,
		fields:   merged,
	}
	logger.sysValues = df
	return logger
}

// Debug writes a write message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use snake_case keys and lower case values if possible.
func Debug(rsFields gcontext.RequestScopedFields, event string, fields ...Fields) {
	defaultLogger.write(rsFields, event, nil, DebugSev, fields...)
}

// Debug writes a write message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Debug(event string, fields ...Fields) {
	logger.write(logger.rsFields, event, nil, DebugSev, fields...)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func Info(rsFields gcontext.RequestScopedFields, event string, fields ...Fields) {
	defaultLogger.write(rsFields, event, nil, InfoSev, fields...)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Info(event string, fields ...Fields) {
	logger.write(logger.rsFields, event, nil, InfoSev, fields...)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func Warn(rsFields gcontext.RequestScopedFields, event string, fields ...Fields) {
	defaultLogger.write(rsFields, event, nil, WarnSev, fields...)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Warn(event string, fields ...Fields) {
	logger.write(logger.rsFields, event, nil, WarnSev, fields...)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func Error(rsFields gcontext.RequestScopedFields, event string, err error, fields ...Fields) {
	defaultLogger.write(rsFields, event, err, ErrorSev, fields...)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Error(event string, err error, fields ...Fields) {
	logger.write(logger.rsFields, event, err, ErrorSev, fields...)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func Fatal(rsFields gcontext.RequestScopedFields, event string, err error, fields ...Fields) {
	event = defaultLogger.write(rsFields, event, err, FatalSev, fields...)

	// time to panic!
	panic(event)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Fatal(event string, err error, fields ...Fields) {
	event = logger.write(logger.rsFields, event, err, FatalSev, fields...)

	// time to panic!
	panic(event)
}

// Event method uses expressive syntax format: logger.Event("event_name").Fields(fields...).Info("message")
func (logger Logger) Event(event string) *Segment {

	return &Segment{
		logger: logger,
		event: event,
		fields: Fields{},
	}
}

func (logger Logger) write(rsFields gcontext.RequestScopedFields, event string, err error, sev string, fields ...Fields) string {
	event = helper.ToSnakeCase(event)

	if sevLevel.shouldLog(sev) {
		system := logger.sysValues.getSystemValues(rsFields, event, sev)
		if err != nil {
			system = logger.sysValues.getErrorValues(err, system)
		}

		properties := logger.fields.Merge(fields...)
		logger.writer.WriteFields(system, properties)
	}

	return event
}


