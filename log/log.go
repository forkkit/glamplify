package log

import (
	"context"
	"io"
	"os"
	"sync"
)

// Config for setting initial values for Logger
type Config struct {
	Output     io.Writer
}

// FieldLogger wraps the standard library logger and add structured types as quoted key value pairs
type FieldLogger struct {
	mutex      sync.Mutex
	output     io.Writer
}

// So that you don't even need to create a new logger
var (
	internal = newLogger(func(conf *Config) {
	})
)

// New creates a new FieldLogger. The optional configure func lets you set values on the underlying standard logger.
// eg. SetOutput
func newLogger(configure ...func(*Config)) *FieldLogger { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

	logger := &FieldLogger{}
	conf := Config{
		Output:     os.Stdout,

	}
	for _, config := range configure {
		config(&conf)
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.output = conf.Output

	return logger
}

// WithScope lets you add types to a scoped logger. Useful for Http Web Request where you want to track user, requestid, etc.
func WithScope(ctx context.Context, fields ...Fields) *Scope {
	return newScope(ctx, internal, fields...)
}

func (logger *FieldLogger) withScope(ctx context.Context, fields ...Fields) *Scope {
	return newScope(ctx, logger, fields...)
}

func (logger *FieldLogger) debug(event string, meta Fields, fields ...Fields) {
	logger.writeFields(event, meta, fields...)
}

func (logger *FieldLogger) info(event string, meta Fields, fields ...Fields) {
	logger.writeFields(event, meta, fields...)
}

func (logger *FieldLogger) warn(event string, meta Fields, fields ...Fields) {
	logger.writeFields(event, meta, fields...)
}

func (logger *FieldLogger) error(event string, meta Fields, fields ...Fields) {
	logger.writeFields(event, meta, fields...)
}

func (logger *FieldLogger) fatal(event string, meta Fields, fields ...Fields) {
	logger.writeFields(event, meta, fields...)

	// time to panic!
	panic(event)
}

func (logger *FieldLogger) writeFields(event string, meta Fields, fields ...Fields) {
	merged := Fields{}
	user := merged.Merge(fields...)
	if len(user) > 0 {
		meta[event] = user
	}
	str := meta.Serialize()
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
