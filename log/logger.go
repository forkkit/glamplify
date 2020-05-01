package log

import (
	"github.com/cultureamp/glamplify/helper"
)

type Config struct {
	TraceId  string
	User     string
	Customer string
}

// Logger
type Logger struct {
	cfg       *Config
	writer    *FieldWriter
	fields    Fields
	defValues *DefaultValues
}

var (
	internal = NewWriter(func(conf *WriterConfig) {})
)

// New creates a *Logger with optional fields. Useful for when you want to add a field to all subsequent logging calls eg. request_id, etc.
func New(cfg Config, fields ...Fields) *Logger {
	return newLogger(cfg, internal, fields...)
}

// Useful for CLI applications that want to write to stderr or file etc.
func NewWitCustomWriter(cfg Config, writer *FieldWriter, fields ...Fields) *Logger {
	return newLogger(cfg, writer, fields...)
}

func newLogger(cfg Config, writer *FieldWriter, fields ...Fields) *Logger {

	df := newDefaultValues()

	merged := Fields{}
	merged = merged.Merge(fields...)
	logger := &Logger{
		cfg:    &cfg,
		writer: writer,
		fields: merged,
	}
	logger.defValues = df
	return logger
}

// Debug writes a debug message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use lower-case keys and values if possible.
func (logger Logger) Debug(event string, fields ...Fields) {
	event = helper.ToSnakeCase(event)
	meta := logger.defValues.getDefaults(logger.cfg, event, DebugSev)
	merged := logger.fields.Merge(fields...)
	logger.writer.debug(event, meta, merged)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (logger Logger) Info(event string, fields ...Fields) {
	event = helper.ToSnakeCase(event)
	meta := logger.defValues.getDefaults(logger.cfg, event, InfoSev)
	merged := logger.fields.Merge(fields...)
	logger.writer.info(event, meta, merged)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (logger Logger) Warn(event string, fields ...Fields) {
	event = helper.ToSnakeCase(event)
	meta := logger.defValues.getDefaults(logger.cfg, event, WarnSev)
	merged := logger.fields.Merge(fields...)
	logger.writer.warn(event, meta, merged)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (logger Logger) Error(err error, fields ...Fields) {
	event := helper.ToSnakeCase(err.Error())
	meta := logger.defValues.getDefaults(logger.cfg, event, ErrorSev)
	meta = logger.defValues.getErrorDefaults(err, meta)
	merged := logger.fields.Merge(fields...)
	logger.writer.error(event, meta, merged)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (logger Logger) Fatal(err error, fields ...Fields) {
	event := helper.ToSnakeCase(err.Error())
	meta := logger.defValues.getDefaults(logger.cfg, event, FatalSev)
	meta = logger.defValues.getErrorDefaults(err, meta)
	merged := logger.fields.Merge(fields...)
	logger.writer.fatal(event, meta, merged)
}
