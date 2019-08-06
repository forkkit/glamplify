package log

import (
	"io"
	"os"
	"sync"
)

// RFC3339Milli is the standard RFC3339 format with added milliseconds
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type Config struct {
	Output io.Writer
	TimeFormat string
}

// FieldLogger wraps the standard library logger and add structured fields as quoted key value pairs
type FieldLogger struct {
	mutex sync.Mutex
	output io.Writer
	timeFormat string
}

// So that you don't even need to create a new logger
var (
	internal = New()
)

// New creates a new FieldLogger. The optional configure func lets you set values on the underlying standard logger.
// eg. SetOutput
func New(configure...func(*Config)) *FieldLogger { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

	logger := &FieldLogger{}
	conf := Config{
		Output:     os.Stdout,
		TimeFormat: RFC3339Milli,
	}
	for _, config := range configure  {
		config(&conf)
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.output = conf.Output
	logger.timeFormat = conf.TimeFormat

	return logger
}

func WithScope(fields Fields) *Scope {
	return newScope(internal, fields)
}

func (logger *FieldLogger)  WithScope(fields Fields) *Scope {
	return newScope(logger, fields)
}

// Debug writes a debug message with optional fields to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {level="debug", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func Debug(message string, fields ...Fields) error {
	return internal.Debug(message, fields...)
}

// Debug writes a debug message with optional fields to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {level="debug", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func (logger FieldLogger) Debug(message string, fields ...Fields) error {
	meta := Fields{
		"host":     hostName(),
		"msg":      message,
		"pid":      processID(),
		"process":  processName(),
		"severity": "DEBUG",
		"time":     timeNow(logger.timeFormat),
	}

	merged := meta.merge(fields...)
	str := merged.serialize()
	return logger.write(str)
}

// Print writes a message with optional fields to the underlying standard logger.
// Useful to normal tracing that should be captured during standard operating behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func Print(message string, fields ...Fields) error {
	return internal.Print(message, fields...)
}

// Print writes a message with optional fields to the underlying standard logger.
// Useful to normal tracing that should be captured during standard operating behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func (logger FieldLogger) Print(message string, fields ...Fields) error {
	meta := Fields{
		"msg":      message,
		"severity": "INFO",
		"time":     timeNow(logger.timeFormat),
	}

	merged := meta.merge(fields...)
	str := merged.serialize()
	return logger.write(str)
}

// Error writes a error message with optional fields to the underlying standard logger.
// Useful to trace errors should be captured always.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {level="error", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func Error(err error, fields ...Fields) error {
	return internal.Error(err, fields...)
}

// Error writes a error message with optional fields to the underlying standard logger.
// Useful to trace errors should be captured always.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {level="error", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func (logger FieldLogger) Error(err error, fields ...Fields) error {
	meta := Fields{
		"arch":		targetArch(),
		"error":    err.Error(),
		"host":     hostName(),
		"os":		targetOS(),
		"pid":      processID(),
		"process":  processName(),
		"severity": "ERROR",
		"time":     timeNow(logger.timeFormat),
	}

	merged := meta.merge(fields...)
	str := merged.serialize()
	return logger.write(str)
}

func (logger *FieldLogger) write(str string) error {

	// Note: Making this faster is a good thing (while we are a sync logger - async logger is a different story)
	// So we don't use the stdlib logger.Print(), but rather have our own optimized version
	// Which does less, but is 3-10x faster

	// alloc a slice to contain the string and possible '\n'
	buffer := make([]byte, len(str) + 1)
	buffer = append(buffer, str...)
	if len(str) == 0 || str[len(str)-1] != '\n' {
		buffer = append(buffer, '\n')
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	_, err := logger.output.Write(buffer)
	return err
}

