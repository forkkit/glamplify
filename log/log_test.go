package log_test

import (
	"bytes"
	"testing"

	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
)

func TestDefaultWarn_Success(t *testing.T) {

	logger := log.Get()

	hook := newMemHook()
	logger.AddHook(hook)

	logger.Warn("test")

	assert.Assert(t, hook.calls == 1, "Hooks calls was: %d. Expected: 1", hook.calls)
	msg := hook.memBuffer.String()
	assert.Assert(t, msg == "test", "Logger was: '%s'. Expected: 'test'", msg)
}

func TestDefaultWarnf_Success(t *testing.T) {

	logger := log.Get()

	hook := newMemHook()
	logger.AddHook(hook)

	logger.Warnf("test %d", 123)

	assert.Assert(t, hook.calls == 1, "Hooks calls was: %d. Expected: 1", hook.calls)
	msg := hook.memBuffer.String()
	assert.Assert(t, msg == "test 123", "Logger was: '%s'. Expected: 'test 123'", msg)
}

func TestDefaultWarnWithFields_Success(t *testing.T) {

	logger := log.Get()

	hook := newMemHook()
	logger.AddHook(hook)

	logger.WarnWithFields(
		log.Fields{
			"string": "hello",
			"int":    123,
			"float":  42.48,
		},
		"test with fields")

	assert.Assert(t, hook.calls == 1, "Hooks calls was: %d. Expected: 1", hook.calls)
	msg := hook.memBuffer.String()
	assert.Assert(t, msg == "test with fields", "Logger was: '%s'. Expected: 'test with fields'", msg)
	assert.Assert(t, hook.entry != nil)
	assert.Assert(t, hook.entry.Fields != nil)

	v, ok := hook.entry.Fields["string"]
	assert.Assert(t, ok)
	assert.Assert(t, v == "hello")

	v, ok = hook.entry.Fields["int"]
	assert.Assert(t, ok)
	assert.Assert(t, v == 123)

	v, ok = hook.entry.Fields["float"]
	assert.Assert(t, ok)
	assert.Assert(t, v == 42.48)
}

func TestDefaultDebug_NotCalled(t *testing.T) {

	logger := log.Get()

	hook := newMemHook()
	logger.AddHook(hook)

	logger.Debug("test")

	assert.Assert(t, hook.calls == 0, "Hooks calls was: %d. Expected: 0", hook.calls)
	msg := hook.memBuffer.String()
	assert.Assert(t, msg == "", "Logger was: '%s'. Expected: ''", msg)
}

func TestNullLogger_DoesNothing(t *testing.T) {
	logger := log.GetFor("something_that_does_not_exist")

	logger.Debug("test")

	assert.Assert(t, logger != nil, "Logger was nil")
	//assert.Assert(t, logger.(*log.nullLogger) != nil, "Not NullLogger")
}

type memHook struct {
	memBuffer *bytes.Buffer
	entry     *log.Entry
	calls     int
}

func newMemHook() *memHook {
	return &memHook{
		memBuffer: &bytes.Buffer{},
		calls:     0,
	}
}

func (h *memHook) Fire(entry *log.Entry) {
	// Write entry.Message to memBuffer
	h.memBuffer.WriteString(entry.Message)
	h.entry = entry
	h.calls++
}
