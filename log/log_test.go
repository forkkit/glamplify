package log_test

import (
	"bytes"
	"errors"
	stdlog "log"
	"os"
	"strings"
	"testing"

	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
)

func TestDebug_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	config := func(stdlogger *stdlog.Logger) {
		stdlogger.SetOutput(memBuffer)
	}
	logger := log.New(config)

	logger.Debug("details")

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "details"), "Logger was: '%s'. Expected: 'details'", msg)
	assert.Assert(t, strings.Contains(msg, "level=\"debug\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, !strings.Contains(msg, "string=\"hello\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, !strings.Contains(msg, "int=\"123\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, !strings.Contains(msg, "float=\"42.48\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
}

func TestDebugWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	config := func(stdlogger *stdlog.Logger) {
		stdlogger.SetOutput(memBuffer)
	}
	logger := log.New(config)

	logger.Debug("details", log.Fields{
		"string": "hello",
		"int":    123,
		"float":  42.48,
	})

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "details"), "Logger was: '%s'. Expected: 'details'", msg)
	assert.Assert(t, strings.Contains(msg, "level=\"debug\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, strings.Contains(msg, "string=\"hello\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, strings.Contains(msg, "int=\"123\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, strings.Contains(msg, "float=\"42.48\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
}

func TestPrint_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	config := func(stdlogger *stdlog.Logger) {
		stdlogger.SetOutput(memBuffer)
	}
	logger := log.New(config)

	logger.Print("info")

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "info"), "Logger was: '%s'. Expected: 'test'", msg)
	assert.Assert(t, !strings.Contains(msg, "string=\"hello\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, !strings.Contains(msg, "int=\"123\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, !strings.Contains(msg, "float=\"42.48\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
}

func TestPrintWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	config := func(stdlogger *stdlog.Logger) {
		stdlogger.SetOutput(memBuffer)
	}
	logger := log.New(config)

	logger.Print("info", log.Fields{
		"string": "hello",
		"int":    123,
		"float":  42.48,
	})

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "info"), "Logger was: '%s'. Expected: 'test'", msg)
	assert.Assert(t, strings.Contains(msg, "string=\"hello\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, strings.Contains(msg, "int=\"123\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, strings.Contains(msg, "float=\"42.48\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
}

func TestError_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	config := func(stdlogger *stdlog.Logger) {
		stdlogger.SetOutput(memBuffer)
	}
	logger := log.New(config)

	logger.Error(errors.New("error"))

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "error"), "Logger was: '%s'. Expected: 'error'", msg)
	assert.Assert(t, strings.Contains(msg, "level=\"error\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, !strings.Contains(msg, "string=\"hello\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, !strings.Contains(msg, "int=\"123\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, !strings.Contains(msg, "float=\"42.48\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
}

func TestErrorWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	config := func(stdlogger *stdlog.Logger) {
		stdlogger.SetOutput(memBuffer)
	}
	logger := log.New(config)

	logger.Error(errors.New("error"), log.Fields{
		"string": "hello",
		"int":    123,
		"float":  42.48,
	})

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "error"), "Logger was: '%s'. Expected: 'error'", msg)
	assert.Assert(t, strings.Contains(msg, "level=\"error\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, strings.Contains(msg, "string=\"hello\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, strings.Contains(msg, "int=\"123\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
	assert.Assert(t, strings.Contains(msg, "float=\"42.48\""), "Logger was: '%s'. Expected: 'level=\"DEBUG\"'", msg)
}

func TestLogSomeRealMessages(t *testing.T) {

	config := func(stdlogger *stdlog.Logger) {
		stdlogger.SetOutput(os.Stderr)
	}
	logger := log.New(config)

	// You should see these printed out, all correctly formatted.
	logger.Debug("details", log.Fields{
		"string": "hello",
		"int":    123,
		"float":  42.48,
	})

	logger.Print("info", log.Fields{
		"string": "hello",
		"int":    123,
		"float":  42.48,
	})
	logger.Error(errors.New("error"), log.Fields{
		"string": "hello",
		"int":    123,
		"float":  42.48,
	})
}
