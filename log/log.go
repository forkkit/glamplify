package log

import (
	"context"
	"fmt"
	"github.com/cultureamp/glamplify/constants"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Config for setting initial values for Logger
type Config struct {
	Output     io.Writer
	TimeFormat string
}

// FieldLogger wraps the standard library logger and add structured types as quoted key value pairs
type FieldLogger struct {
	mutex      sync.Mutex
	output     io.Writer
	timeFormat string
	defValues  *DefaultValues
}

// So that you don't even need to create a new logger
var (
	internal = New(func(conf *Config) {
	})
)

// New creates a new FieldLogger. The optional configure func lets you set values on the underlying standard logger.
// eg. SetOutput
func New(configure ...func(*Config)) *FieldLogger { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

	logger := &FieldLogger{}
	conf := Config{
		Output:     os.Stdout,
		TimeFormat: constants.RFC3339Milli,
	}
	for _, config := range configure {
		config(&conf)
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.output = conf.Output
	logger.timeFormat = conf.TimeFormat
	logger.defValues = NewDefaultValues(logger.timeFormat)

	return logger
}

// WithScope lets you add types to a scoped logger. Useful for Http Web Request where you want to track user, requestid, etc.
func WithScope(fields Fields) *Scope {
	return newScope(internal, fields)
}

// WithScope lets you add types to a scoped logger. Useful for Http Web Request where you want to track user, requestid, etc.
func (logger *FieldLogger) WithScope(fields Fields) *Scope {
	return newScope(logger, fields)
}

// Debug writes a debug message with optional types to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use lower-case keys and values if possible.
func Debug(message string, fields ...Fields) {
	internal.Debug(message, fields...)
}

// Debug writes a debug message with optional types to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use lower-case keys and values if possible.
func (logger *FieldLogger) Debug(message string, fields ...Fields) {
	meta := logger.defValues.GetDefaults(message, constants.DebugSevLogValue)
	logger.writeFields(meta, fields...)
}

// Info writes a message with optional types to the underlying standard logger.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func Info(message string, fields ...Fields) {
	internal.Info(message, fields...)
}

// Info writes a message with optional types to the underlying standard logger.
// Useful form normal tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (logger *FieldLogger) Info(message string, fields ...Fields) {
	meta := logger.defValues.GetDefaults(message, constants.InfoSevLogValue)
	logger.writeFields(meta, fields...)
}

// Warn writes a message with optional types to the underlying standard logger.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func Warn(message string, fields ...Fields) {
	internal.Warn(message, fields...)
}

// Warn writes a message with optional types to the underlying standard logger.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (logger *FieldLogger) Warn(message string, fields ...Fields) {
	meta := logger.defValues.GetDefaults(message, constants.WarnSevLogValue)
	logger.writeFields(meta, fields...)
}

// Error writes a error message with optional types to the underlying standard logger.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func Error(err error, fields ...Fields) {
	internal.Error(err, fields...)
}

// Error writes a error message with optional types to the underlying standard logger.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (logger *FieldLogger) Error(err error, fields ...Fields) {
	meta := logger.defValues.GetDefaults(strings.TrimSpace(err.Error()), constants.ErrorSevLogValue)
	meta = logger.defValues.GetErrorDefaults(err, meta)
	logger.writeFields(meta, fields...)
}

// Fatal writes a error message with optional types to the underlying standard logger and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func Fatal(err error, fields ...Fields) {
	internal.Fatal(err, fields...)
}

// Fatal writes a error message with optional types to the underlying standard logger and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (logger *FieldLogger) Fatal(err error, fields ...Fields) {
	meta := logger.defValues.GetDefaults(strings.TrimSpace(err.Error()), constants.FatalSevLogValue)
	meta = logger.defValues.GetErrorDefaults(err, meta)
	logger.writeFields(meta, fields...)

	// time to panic!
	panic(strings.TrimSpace(err.Error()))
}

func Audit(ctx context.Context, event string, err error, fields ...Fields) {
	internal.Audit(ctx, event, err, fields...)
}

func (logger *FieldLogger) Audit(ctx context.Context, event string, err error, fields ...Fields) {
	meta := logger.defValues.GetDefaults(constants.EmptyString, constants.AuditSevLogValue)
	if err != nil {
		meta = logger.defValues.GetErrorDefaults(err, meta)
	}
	meta = logger.defValues.GetAuditDefaults(ctx, event, meta)
	logger.writeFields(meta, fields...)
}


func DurationAsISO8601(duration time.Duration) string {
	return internal.DurationAsISO8601(duration)
}

func (logger FieldLogger) DurationAsISO8601(duration time.Duration) string {
	return fmt.Sprintf("P%gS", duration.Seconds())
}


func (logger *FieldLogger) writeFields(meta Fields, fields ...Fields) {
	merged := meta.Merge(fields...)
	str := merged.Serialize()
	logger.write(str)
}

func (logger *FieldLogger) write(str string) {

	// Note: Making this faster is a good thing (while we are a sync logger - async logger is a different story)
	// So we don't use the stdlib logger.Print(), but rather have our own optimized version
	// Which does less, but is 3-10x faster

	// alloc a slice to contain the string and possible '\n'
	length := len(str)
	buffer := make([]byte, length+1)
	copy(buffer[:], str)
	if len(str) == 0 || str[length-1] != '\n' {
		copy(buffer[length:], "\n")
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	// This can return an error, but we just swallow it here as what can we or a client really do? Try and log it? :)
	logger.output.Write(buffer)
}
