package log

import (
	"bytes"
	"testing"

	"gotest.tools/assert"
)

// Testify?

// TestRunner todo...
func TestDefaultWarn(t *testing.T) {

	logger := LoggerFactory.Get("default")

	hook := newMemHook()
	logger.AddHook(hook)

	logger.Warn("test")

	assert.Assert(t, hook.calls == 1, "Hooks calls was: %d. Expected: 1", hook.calls)
	msg := hook.memBuffer.String()
	assert.Assert(t, msg == "test", "Logger was: '%s'. Expected: 'test'", msg)
}

func TestDefaultDebug(t *testing.T) {

	logger := LoggerFactory.Get("default")

	hook := newMemHook()
	logger.AddHook(hook)

	logger.Debug("test")

	assert.Assert(t, hook.calls == 0, "Hooks calls was: %d. Expected: 0", hook.calls)
	msg := hook.memBuffer.String()
	assert.Assert(t, msg == "", "Logger was: '%s'. Expected: ''", msg)
}

func TestNullLogger(t *testing.T) {
	logger := LoggerFactory.Get("something_that_does_not_exist")

	logger.Debug("test")

	assert.Assert(t, logger != nil, "Logger was nil")
	assert.Assert(t, logger.(*nullLogger) != nil, "Not NullLogger")
}

type memHook struct {
	memBuffer *bytes.Buffer
	calls     int
}

func newMemHook() *memHook {
	return &memHook{
		memBuffer: &bytes.Buffer{},
		calls:     0,
	}
}

func (h *memHook) Fire(entry *LogEntry) {
	// Write entry.Message to memBuffer
	h.memBuffer.WriteString(entry.Message)
	h.calls++
}
