package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RFC3339Milli is the standard RFC3339 format with added milliseconds
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

// Fields type, used to pass to Debug, Print and Error.
type Fields map[string]interface{}

// FieldLogger wraps the standard library logger and add structured fields as quoted key value pairs
type FieldLogger struct {
	outputMutex sync.Mutex
	output io.Writer

	timeFormat string
	context    Fields
}

// So that you don't even need to create a new logger
var internal = New()

// New creates a new FieldLogger. The optional configure func lets you set values on the underlying standard logger.
// eg. SetOutput
func New() *FieldLogger { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

	logger := &FieldLogger{}
	logger.output = os.Stdout
	logger.timeFormat = RFC3339Milli
	logger.context = Fields{}

	return logger
}

// AddContext will add a key-value pair to every logging message
func (logger *FieldLogger) AddContext(key string, value interface{}) {
	logger.context[key] = value
}

// SetOutput set the output to anything that supports io.Writer
func (logger *FieldLogger) SetOutput(writer io.Writer) {
	logger.outputMutex.Lock()
	defer logger.outputMutex.Unlock()

	logger.output = writer
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

	str := combine(meta, logger.context, fields...)
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

	str := combine(meta, logger.context, fields...)
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

	str := combine(meta, logger.context, fields...)
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

	logger.outputMutex.Lock()
	defer logger.outputMutex.Unlock()

	_, err := logger.output.Write(buffer)
	return err
}

func combine(meta Fields, context Fields, fields ...Fields) string {

	var str []string

	count, pre := serialize(meta)
	if count > 0 {
		str = append(str, pre...)
	}

	count, ctx := serialize(context)
	if count > 0 {
		str = append(str, ctx...)
	}

	for _, f := range fields {
		count, post := serialize(f)
		if count > 0 {
			str = append(str, post...)
		}
	}

	sort.Strings(str)
	return strings.Join(str, " ")
}

func serialize(fields Fields) (int, []string) {
	var pairs []string
	for k, v := range fields {
		vs := fmt.Sprintf("%v", v)

		pairs = append(pairs, quoteIfRequired(k)+"="+quoteIfRequired(vs))
	}

	return len(pairs), pairs
}

func quoteIfRequired(input string) string {
	if strings.Contains(input, " ") {
		input = strconv.Quote(input)
	}
	return input
}

func timeNow(format string) string {
	return time.Now().Format(format)
}

var host string
var hostOnce sync.Once

func hostName() string {

	var err error
	hostOnce.Do(func() {
		host, err = os.Hostname()
		if err != nil {
			host = "<unknown>"
		}
	})

	return host
}

func processName() string {
	name := os.Args[0]
	if len(name) > 0 {
		name = filepath.Base(name)
	}

	return name
}

var pid int
var pidOnce sync.Once

func processID() int {
	pidOnce.Do(func() {
		pid = os.Getpid()
	})

	return pid
}

func targetArch() string {
	return runtime.GOARCH
}

func targetOS() string {
	return runtime.GOOS
}
