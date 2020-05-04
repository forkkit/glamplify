package log

import (
	"github.com/cultureamp/glamplify/helper"
)

// Logger
type Logger struct {
	tFields   TransactionFields
	writer    *FieldWriter
	fields    Fields
	defValues *DefaultValues
}

var (
	internalWriter = NewWriter(func(conf *WriterConfig) {})
	defaultLogger  = NewWitCustomWriter(TransactionFields{}, internalWriter)
)

// New creates a *Logger with optional fields. Useful for when you want to add a field to all subsequent logging calls eg. request_id, etc.
func New(tFields TransactionFields, fields ...Fields) *Logger {
	return newLogger(tFields, internalWriter, fields...)
}

// Useful for CLI applications that want to write to stderr or file etc.
func NewWitCustomWriter(tFields TransactionFields, writer *FieldWriter, fields ...Fields) *Logger {
	return newLogger(tFields, writer, fields...)
}

func newLogger(tFields TransactionFields, writer *FieldWriter, fields ...Fields) *Logger {

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
func Debug(tFields TransactionFields, event string, fields ...Fields) {
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
func Info(tFields TransactionFields, event string, fields ...Fields) {
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
func Warn(tFields TransactionFields, event string, fields ...Fields) {
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
func Error(tFields TransactionFields, err error, fields ...Fields) {
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
func Fatal(tFields TransactionFields, err error, fields ...Fields) {
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

func (logger Logger) write(tFields TransactionFields, event string, sev string, fields ...Fields) string {
	event = helper.ToSnakeCase(event)
	meta := logger.defValues.getDefaults(tFields, event, sev)
	merged := logger.fields.Merge(fields...)
	logger.writer.writeFields(event, meta, merged)

	return event
}

func (logger Logger) writeError(tFields TransactionFields, err error, sev string, fields ...Fields) string {
	event := helper.ToSnakeCase(err.Error())
	meta := logger.defValues.getDefaults(tFields, event, sev)
	meta = logger.defValues.getErrorDefaults(err, meta)
	merged := logger.fields.Merge(fields...)
	logger.writer.writeFields(event, meta, merged)

	return event
}
