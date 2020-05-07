package log

import (
	"io"
	"os"
	"sync"
)

// WriterConfig for setting initial values for Logger
type WriterConfig struct {
	Output io.Writer
}

// FieldWriter wraps the standard library writer and add structured types as quoted key value pairs
type FieldWriter struct {
	mutex      sync.Mutex
	output     io.Writer
}

// NewWriter creates a new FieldWriter. The optional configure func lets you set values on the underlying standard writer.
// Useful for CLI apps that want to direct logging to a file or stderr
// eg. SetOutput
func NewWriter(configure ...func(*WriterConfig)) *FieldWriter { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

	logger := &FieldWriter{}
	conf := WriterConfig{
		Output: os.Stdout,
	}
	for _, config := range configure {
		config(&conf)
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.output = conf.Output

	return logger
}

func (writer *FieldWriter) writeFields(system Fields, fields ...Fields) {
	merged := Fields{}
	properties := merged.Merge(fields...)
	if len(properties) > 0 {
		system[Properties] = properties
	}
	str := system.ToSnakeCase().serialize()
	writer.write(str)
}

func (writer *FieldWriter) write(str string) {

	// Note: Making this faster is a good thing (while we are a sync writer - async writer is a different story)
	// So we don't use the stdlib writer.Print(), but rather have our own optimized version
	// Which does less, but is 3-10x faster

	// alloc a slice to contain the string and possible '\n'
	length := len(str)
	buffer := make([]byte, length+1)
	copy(buffer[:], str)
	if len(str) == 0 || str[length-1] != '\n' {
		copy(buffer[length:], "\n")
	}

	writer.mutex.Lock()
	defer writer.mutex.Unlock()

	// This can return an error, but we just swallow it here as what can we or a client really do? Try and log it? :)
	writer.output.Write(buffer)
}
