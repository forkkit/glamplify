package log

import (
	"context"
	"github.com/cultureamp/glamplify/helper"
	"net/http"
)

// Logger
type Logger struct {
	rsFields  RequestScopedFields
	fields    Fields
	sysValues *SystemValues
	writer    *FieldWriter
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

// NewFromCtx creates a new logger from a context, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromCtx(ctx context.Context, fields ...Fields) *Logger {

	rsFields, _ := GetRequestScopedFieldsFromCtx(ctx)
	logger := New(rsFields, fields...)
	return logger
}

// NewFromRequest creates a new logger from a http.Request, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromRequest(r *http.Request, fields ...Fields) *Logger {
	logger := NewFromCtx(r.Context(), fields...)
	return logger
}

func newLogger(rsFields RequestScopedFields, writer *FieldWriter, fields ...Fields) *Logger {

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
func Debug(tFields RequestScopedFields, event string, fields ...Fields) {
	defaultLogger.write(tFields, event, DebugSev, fields...)
}

// Debug writes a write message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Debug(event string, fields ...Fields) {
	logger.write(logger.rsFields, event, DebugSev, fields...)
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
	logger.write(logger.rsFields, event, InfoSev, fields...)
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
	logger.write(logger.rsFields, event, WarnSev, fields...)
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
	logger.writeError(logger.rsFields, err, ErrorSev, fields...)
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
	event := logger.writeError(logger.rsFields, err, FatalSev, fields...)

	// time to panic!
	panic(event)
}

func (logger Logger) write(rsFields RequestScopedFields, event string, sev string, fields ...Fields) string {
	event = helper.ToSnakeCase(event)

	system := logger.sysValues.getSystemValues(rsFields, event, sev)

	properties := logger.fields.Merge(fields...)
	logger.writer.writeFields(system, properties)

	return event
}

func (logger Logger) writeError(rsFields RequestScopedFields, err error, sev string, fields ...Fields) string {
	event := helper.ToSnakeCase(err.Error())

	system := logger.sysValues.getSystemValues(rsFields, event, sev)
	systemWithErrors := logger.sysValues.getErrorValues(err, system)

	properties := logger.fields.Merge(fields...)
	logger.writer.writeFields(systemWithErrors, properties)

	return event
}
