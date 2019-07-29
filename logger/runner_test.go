package logger

import (
	"bytes"
	"strings"
	"testing"

	"gotest.tools/assert"
)

// TestRunner todo...
func TestDefaultWarn(t *testing.T) {

	var memBuffer bytes.Buffer

	logger := LoggerFactory.Get("default")
	logger.SetOutput(&memBuffer)

	logger.Warn("test")

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=test"), "Logger was: '%s'. Expected: 'test'", msg)
}

func TestDefaultDebug(t *testing.T) {

	var memBuffer bytes.Buffer

	logger := LoggerFactory.Get("default")
	logger.SetOutput(&memBuffer)

	logger.Debug("test")

	msg := memBuffer.String()
	assert.Assert(t, msg == "", "Logger was: '%s'. Expected: ''", msg)
}

func TestNullLogger(t *testing.T) {
	logger := LoggerFactory.Get("something_that_does_not_exist")

	logger.Debug("test")

	assert.Assert(t, logger != nil, "Logger was nil")
	assert.Assert(t, logger.(*nullLogger) != nil, "Not NullLogger")
}
