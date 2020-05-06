package log

import (
	"context"
	"github.com/cultureamp/glamplify/helper"
	"net/http"
)

// Logger
type Logger struct {
	tFields   RequestScopedFields
	writer    *FieldWriter
	fields    Fields
	defValues *DefaultValues
}

var (
	internalWriter = NewWriter(func(conf *WriterConfig) {})
	defaultLogger  = NewWitCustomWriter(RequestScopedFields{}, internalWriter)
)

// New creates a *Logger with optional fields. Useful for when you want to add a field to all subsequent logging calls eg. request_id, etc.
func New(rsFields RequestScopedFields, fields ...Fields) *Logger {
	return newLogger(rsFields, internalWriter, fields...)
}

// Useful for CLI applications that want to write to stderr or file etc.
func NewWitCustomWriter(rsFields RequestScopedFields, writer *FieldWriter, fields ...Fields) *Logger {
	return newLogger(rsFields, writer, fields...)
}

// NewWithCtx creates a new logger from a context. The context should contain TraceID, CustomerID, UserID, but
// if not then this method will create and add TraceID to the context and return that new context
// which should then be used in further methods (so we have consistent trace_ids)
// To add CustomerID, UserID use ctx := log.AddCustomer(ctx, customerID), ctx := log.AddUser(ctx, userID)
// before calling this method. You can use the jwt helper in the package to get these values from the JWT
func NewWithCtx(ctx context.Context,  fields ...Fields) (context.Context, *Logger) {
	rsFields := NewRequestScopeFieldsFromCtx(ctx)
	ctx = rsFields.AddToCtx(ctx)
	logger := New(rsFields, fields...)
	return ctx, logger
}

// NewWithRequest creates a new logger from a http.Request. The context within the request should contain
// TraceID, CustomerID, UserID, but if not then this method wil create and add TraceID to the context and return
// the http.Request with that new context. This returned http.Request should be used in further methods (so we have
// consistent trace_ids) To add CustomerID, UserID use ctx := log.AddCustomer(ctx, customerID), ctx := log.AddUser(ctx, userID)
// before calling this method. You can use the jwt helper in the package to get these values from the JWT
func NewWithRequest(r *http.Request, fields ...Fields) (*http.Request, *Logger){
	ctx, logger := NewWithCtx(r.Context(), fields...)
	return r.WithContext(ctx), logger
}

func newLogger(tFields RequestScopedFields, writer *FieldWriter, fields ...Fields) *Logger {

	df := newDefaultValues()

	merged := Fields{}
	merged = merged.Merge(fields...)
	logger := &Logger{
		tFields: tFields,
		writer:  writer,
		fields:  merged,
	}
	logger.defValues = df
	return logger
}

// Debug writes a write message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use snake_case keys and lower case values if possible.
func Debug(tFields RequestScopedFields, event string, fields ...Fields) {
	defaultLogger.write(tFields, event, DebugSev, fields...)
}

// Debug writes a write message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Debug(event string, fields ...Fields) {
	logger.write(logger.tFields, event, DebugSev, fields...)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func Info(tFields RequestScopedFields, event string, fields ...Fields) {
	defaultLogger.write(tFields, event, InfoSev, fields...)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Info(event string, fields ...Fields) {
	logger.write(logger.tFields, event, InfoSev, fields...)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func Warn(tFields RequestScopedFields, event string, fields ...Fields) {
	defaultLogger.write(tFields, event, WarnSev, fields...)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Warn(event string, fields ...Fields) {
	logger.write(logger.tFields, event, WarnSev, fields...)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func Error(tFields RequestScopedFields, err error, fields ...Fields) {
	defaultLogger.writeError(tFields, err, ErrorSev, fields...)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Error(err error, fields ...Fields) {
	logger.writeError(logger.tFields, err, ErrorSev, fields...)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func Fatal(tFields RequestScopedFields, err error, fields ...Fields) {
	event := defaultLogger.writeError(tFields, err, FatalSev, fields...)

	// time to panic!
	panic(event)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Fatal(err error, fields ...Fields) {
	event := logger.writeError(logger.tFields, err, FatalSev, fields...)

	// time to panic!
	panic(event)
}

func (logger Logger) write(tFields RequestScopedFields, event string, sev string, fields ...Fields) string {
	event = helper.ToSnakeCase(event)
	meta := logger.defValues.getDefaults(tFields, event, sev)
	merged := logger.fields.Merge(fields...)
	logger.writer.writeFields(event, meta, merged)

	return event
}

func (logger Logger) writeError(tFields RequestScopedFields, err error, sev string, fields ...Fields) string {
	event := helper.ToSnakeCase(err.Error())
	meta := logger.defValues.getDefaults(tFields, event, sev)
	meta = logger.defValues.getErrorDefaults(err, meta)
	merged := logger.fields.Merge(fields...)
	logger.writer.writeFields(event, meta, merged)

	return event
}
