package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// RFC3339Milli is the standard RFC3339 format with added milliseconds
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

// Fields type, used to pass to Debug, Print and Error.
type Fields map[string]interface{}

// FieldLogger wraps the standard library logger and add structured fields as quoted key value pairs
type FieldLogger struct {
	stdLogger  *log.Logger
	timeFormat string
}

// So that you don't even need to create a new logger
var internal = New()

// New creates a new FieldLogger. The optional configure func lets you set values on the underlying standard logger.
// eg. SetOutput
func New() *FieldLogger { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

	logger := &FieldLogger{}
	logger.stdLogger = log.New(os.Stdout, "", 0)
	logger.timeFormat = RFC3339Milli

	return logger
}

// SetOutput set the output to anything that supports io.Writer
func (logger *FieldLogger) SetOutput(writer io.Writer) {
	logger.stdLogger.SetOutput(writer)
}

// SetTimeFormat allows you to change the default time format from "2006-01-02T15:04:05.000Z07:00" to whatever you like
func (logger *FieldLogger) SetTimeFormat(format string) {
	logger.timeFormat = format
}

// Debug writes a debug message with optional fields to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {level="debug", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func Debug(message string, fields ...Fields) {
	internal.Debug(message, fields...)
}

// Debug writes a debug message with optional fields to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {level="debug", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func (logger FieldLogger) Debug(message string, fields ...Fields) {
	meta := Fields{
		"level":   "debug",
		"host":    hostName(),
		"pid":     processID(),
		"process": processName(),
		"time":    timeNow(logger.timeFormat),
	}

	str := combine(meta, message, fields...)
	logger.stdLogger.Print(str)
}

// Print writes a message with optional fields to the underlying standard logger.
// Useful to normal tracing that should be captured during standard operating behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func Print(message string, fields ...Fields) {
	internal.Print(message, fields...)
}

// Print writes a message with optional fields to the underlying standard logger.
// Useful to normal tracing that should be captured during standard operating behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func (logger FieldLogger) Print(message string, fields ...Fields) {
	meta := Fields{
		"time": timeNow(logger.timeFormat),
	}

	str := combine(meta, message, fields...)
	logger.stdLogger.Print(str)
}

// Error writes a error message with optional fields to the underlying standard logger.
// Useful to trace errors should be captured always.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {level="error", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func Error(err error, fields ...Fields) {
	internal.Error(err, fields...)
}

// Error writes a error message with optional fields to the underlying standard logger.
// Useful to trace errors should be captured always.
// All field values will be automatically quoted (keys will not be).
// Debug adds fields {level="error", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "fields message"
// Use lower-case keys and values if possible.
func (logger FieldLogger) Error(err error, fields ...Fields) {
	meta := Fields{
		"level":   "error",
		"host":    hostName(),
		"pid":     processID(),
		"process": processName(),
		"time":    timeNow(logger.timeFormat),
	}

	str := combine(meta, err.Error(), fields...)
	logger.stdLogger.Print(str)
}

func combine(meta Fields, message string, fields ...Fields) string {

	var str []string

	count, pre := serialize(meta)
	if count > 0 {
		str = append(str, pre)
	}

	for _, f := range fields {
		count, post := serialize(f)
		if count > 0 {
			str = append(str, post)
		}
	}

	if len(message) > 0 {
		str = append(str, message)
	}

	return strings.Join(str, " ")
}

func serialize(fields Fields) (int, string) {
	var pairs []string
	for k, v := range fields {
		vs := fmt.Sprintf("%v", v)
		pairs = append(pairs, k+"="+strconv.Quote(vs))
	}
	sort.Strings(pairs)
	return len(pairs), strings.Join(pairs, " ")
}

func timeNow(format string) string {
	return time.Now().Format(format)
}

func hostName() string {
	name, err := os.Hostname()
	if err != nil {
		name = "<unknown>"
	}

	return name
}

func processName() string {
	name := os.Args[0]
	if len(name) > 0 {
		name = filepath.Base(name)
	}

	return name
}

func processID() int {
	return os.Getpid()
}
