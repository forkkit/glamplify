package log

import (
	"github.com/cultureamp/glamplify/helper"
)

// Logger
type Logger struct {
	mFields   MandatoryFields
	writer    *FieldWriter
	fields    Fields
	defValues *DefaultValues
}

var (
	internalWriter = NewWriter(func(conf *WriterConfig) {})
	defaultLogger  = NewWitCustomWriter(MandatoryFields{}, internalWriter)
)

// New creates a *Logger with optional fields. Useful for when you want to add a field to all subsequent logging calls eg. request_id, etc.
func New(mFields MandatoryFields, fields ...Fields) *Logger {
	return newLogger(mFields, internalWriter, fields...)
}

// Useful for CLI applications that want to write to stderr or file etc.
func NewWitCustomWriter(mFields MandatoryFields, writer *FieldWriter, fields ...Fields) *Logger {
	return newLogger(mFields, writer, fields...)
}

func newLogger(mFields MandatoryFields, writer *FieldWriter, fields ...Fields) *Logger {

	df := newDefaultValues()

	merged := Fields{}
	merged = merged.Merge(fields...)
	logger := &Logger{
		mFields: mFields,
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
func Debug(mFields MandatoryFields, event string, fields ...Fields) {
	defaultLogger.write(mFields, event, DebugSev, fields...)
}

// Debug writes a write message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Debug(event string, fields ...Fields) {
	logger.write(logger.mFields, event, DebugSev, fields...)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func Info(mFields MandatoryFields, event string, fields ...Fields) {
	defaultLogger.write(mFields, event, InfoSev, fields...)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Info(event string, fields ...Fields) {
	logger.write(logger.mFields, event, InfoSev, fields...)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func Warn(mFields MandatoryFields, event string, fields ...Fields) {
	defaultLogger.write(mFields, event, WarnSev, fields...)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Warn(event string, fields ...Fields) {
	logger.write(logger.mFields, event, WarnSev, fields...)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func Error(mFields MandatoryFields, err error, fields ...Fields) {
	defaultLogger.writeError(mFields, err, ErrorSev, fields...)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Error(err error, fields ...Fields) {
	logger.writeError(logger.mFields, err, ErrorSev, fields...)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func Fatal(mFields MandatoryFields, err error, fields ...Fields) {
	event := defaultLogger.writeError(mFields, err, FatalSev, fields...)

	// time to panic!
	panic(event)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Fatal(err error, fields ...Fields) {
	event := logger.writeError(logger.mFields, err, FatalSev, fields...)

	// time to panic!
	panic(event)
}

func (logger Logger) write(mFields MandatoryFields, event string, sev string, fields ...Fields) string {
	event = helper.ToSnakeCase(event)
	meta := logger.defValues.getDefaults(mFields, event, sev)
	merged := logger.fields.Merge(fields...)
	logger.writer.writeFields(event, meta, merged)

	return event
}

func (logger Logger) writeError(mFields MandatoryFields, err error, sev string, fields ...Fields) string {
	event := helper.ToSnakeCase(err.Error())
	meta := logger.defValues.getDefaults(mFields, event, sev)
	meta = logger.defValues.getErrorDefaults(err, meta)
	merged := logger.fields.Merge(fields...)
	logger.writer.writeFields(event, meta, merged)

	return event
}
