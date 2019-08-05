package log_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
)

func TestDebug_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New()
	logger.SetOutput(memBuffer)

	logger.Debug("details")

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=details"), "Logger was: '%s'. Expected: 'msg=details'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=DEBUG"), "Logger was: '%s'. Expected: 'level=DEBUG'", msg)
}

func TestDebugWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New()
	logger.SetOutput(memBuffer)

	logger.Debug("details", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=details"), "Logger was: '%s'. Expected: 'msg=details'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=DEBUG"), "Logger was: '%s'. Expected: 'level=DEBUG'", msg)
	assert.Assert(t, strings.Contains(msg, "string=hello"), "Logger was: '%s'. Expected: 'string=hello'", msg)
	assert.Assert(t, strings.Contains(msg, "int=123"), "Logger was: '%s'. Expected: 'int=123'", msg)
	assert.Assert(t, strings.Contains(msg, "float=42.48"), "Logger was: '%s'. Expected: 'float=42.48'", msg)
	assert.Assert(t, strings.Contains(msg, "string2=\"hello world\""), "Logger was: '%s'. Expected: 'string2=\"hello world\"'", msg)
	assert.Assert(t, strings.Contains(msg, "\"string3 space\"=world"), "Logger was: '%s'. Expected: '\"string3 space\"=world'", msg)
}

func TestPrint_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New()
	logger.SetOutput(memBuffer)

	logger.Print("info")

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=info"), "Logger was: '%s'. Expected: 'msg=info'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=INFO"), "Logger was: '%s'. Expected: 'level=DEBUG'", msg)
}

func TestPrintWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New()
	logger.SetOutput(memBuffer)

	logger.Print("info", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=info"), "Logger was: '%s'. Expected: 'msg=info'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=INFO"), "Logger was: '%s'. Expected: 'level=DEBUG'", msg)
	assert.Assert(t, strings.Contains(msg, "string=hello"), "Logger was: '%s'. Expected: 'string=hello'", msg)
	assert.Assert(t, strings.Contains(msg, "int=123"), "Logger was: '%s'. Expected: 'int=123'", msg)
	assert.Assert(t, strings.Contains(msg, "float=42.48"), "Logger was: '%s'. Expected: 'float=42.48'", msg)
	assert.Assert(t, strings.Contains(msg, "string2=\"hello world\""), "Logger was: '%s'. Expected: 'string2=\"hello world\"'", msg)
	assert.Assert(t, strings.Contains(msg, "\"string3 space\"=world"), "Logger was: '%s'. Expected: '\"string3 space\"=world'", msg)
}

func TestError_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New()
	logger.SetOutput(memBuffer)

	logger.Error(errors.New("error"))

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "error=error"), "Logger was: '%s'. Expected: 'error=error'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=ERROR"), "Logger was: '%s'. Expected: 'severity=ERROR'", msg)
}

func TestErrorWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New()
	logger.SetOutput(memBuffer)

	logger.Error(errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "error=error"), "Logger was: '%s'. Expected: 'error=error'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=ERROR"), "Logger was: '%s'. Expected: 'severity=ERROR'", msg)
	assert.Assert(t, strings.Contains(msg, "string=hello"), "Logger was: '%s'. Expected: 'string=hello'", msg)
	assert.Assert(t, strings.Contains(msg, "int=123"), "Logger was: '%s'. Expected: 'int=123'", msg)
	assert.Assert(t, strings.Contains(msg, "float=42.48"), "Logger was: '%s'. Expected: 'float=2.48'", msg)
	assert.Assert(t, strings.Contains(msg, "string2=\"hello world\""), "Logger was: '%s'. Expected: 'string2=\"hello world\"'", msg)
	assert.Assert(t, strings.Contains(msg, "\"string3 space\"=world"), "Logger was: '%s'. Expected: '\"string3 space\"=world'", msg)
}

func TestLogSomeRealMessages(t *testing.T) {

	logger := log.New()
	logger.SetOutput(os.Stderr)
	logger.AddContext("app", "mytest-app.exe")

	// You should see these printed out, all correctly formatted.
	logger.Debug("details", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Print("info", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	logger.Error(errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
}
