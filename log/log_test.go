package log_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
)

func TestDebug_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	err := logger.Debug("details")
	assert.Assert(t, err == nil)

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=details"), "Logger was: '%s'. Expected: 'msg=details'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=DEBUG"), "Logger was: '%s'. Expected: 'severity=DEBUG'", msg)
}

func TestDebugWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	err := logger.Debug("details", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	assert.Assert(t, err == nil)

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=details"), "Logger was: '%s'. Expected: 'msg=details'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=DEBUG"), "Logger was: '%s'. Expected: 'severity=DEBUG'", msg)
	assert.Assert(t, strings.Contains(msg, "string=hello"), "Logger was: '%s'. Expected: 'string=hello'", msg)
	assert.Assert(t, strings.Contains(msg, "int=123"), "Logger was: '%s'. Expected: 'int=123'", msg)
	assert.Assert(t, strings.Contains(msg, "float=42.48"), "Logger was: '%s'. Expected: 'float=42.48'", msg)
	assert.Assert(t, strings.Contains(msg, "string2=\"hello world\""), "Logger was: '%s'. Expected: 'string2=\"hello world\"'", msg)
	assert.Assert(t, strings.Contains(msg, "\"string3 space\"=world"), "Logger was: '%s'. Expected: '\"string3 space\"=world'", msg)
}

func TestPrint_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	err := logger.Print("info")
	assert.Assert(t, err == nil)

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=info"), "Logger was: '%s'. Expected: 'msg=info'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=INFO"), "Logger was: '%s'. Expected: 'severity=INFO'", msg)
}

func TestPrintWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	err := logger.Print("info", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	assert.Assert(t, err == nil)

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=info"), "Logger was: '%s'. Expected: 'msg=info'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=INFO"), "Logger was: '%s'. Expected: 'severity=INFO'", msg)
	assert.Assert(t, strings.Contains(msg, "string=hello"), "Logger was: '%s'. Expected: 'string=hello'", msg)
	assert.Assert(t, strings.Contains(msg, "int=123"), "Logger was: '%s'. Expected: 'int=123'", msg)
	assert.Assert(t, strings.Contains(msg, "float=42.48"), "Logger was: '%s'. Expected: 'float=42.48'", msg)
	assert.Assert(t, strings.Contains(msg, "string2=\"hello world\""), "Logger was: '%s'. Expected: 'string2=\"hello world\"'", msg)
	assert.Assert(t, strings.Contains(msg, "\"string3 space\"=world"), "Logger was: '%s'. Expected: '\"string3 space\"=world'", msg)
}

func TestError_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})


	err := logger.Error(errors.New("error"))
	assert.Assert(t, err == nil)

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "error=error"), "Logger was: '%s'. Expected: 'error=error'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=ERROR"), "Logger was: '%s'. Expected: 'severity=ERROR'", msg)
}

func TestErrorWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	err := logger.Error(errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	assert.Assert(t, err == nil)

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "error=error"), "Logger was: '%s'. Expected: 'error=error'", msg)
	assert.Assert(t, strings.Contains(msg, "severity=ERROR"), "Logger was: '%s'. Expected: 'severity=ERROR'", msg)
	assert.Assert(t, strings.Contains(msg, "string=hello"), "Logger was: '%s'. Expected: 'string=hello'", msg)
	assert.Assert(t, strings.Contains(msg, "int=123"), "Logger was: '%s'. Expected: 'int=123'", msg)
	assert.Assert(t, strings.Contains(msg, "float=42.48"), "Logger was: '%s'. Expected: 'float=2.48'", msg)
	assert.Assert(t, strings.Contains(msg, "string2=\"hello world\""), "Logger was: '%s'. Expected: 'string2=\"hello world\"'", msg)
	assert.Assert(t, strings.Contains(msg, "\"string3 space\"=world"), "Logger was: '%s'. Expected: '\"string3 space\"=world'", msg)
}

func TestScope(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	scope := logger.WithScope(log.Fields{
		"requestID": 123,
	})

	err := scope.Debug("details")
	assert.Assert(t, err == nil)

	msg := memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=details"), "Logger was: '%s'. Expected: 'msg=details'", msg)
	assert.Assert(t, strings.Contains(msg, "requestID=123"), "Logger was: '%s'. Expected: 'requestID=123'", msg)

	memBuffer.Reset()

	err = scope.Print("info")
	assert.Assert(t, err == nil)

	msg = memBuffer.String()
	assert.Assert(t, strings.Contains(msg, "msg=info"), "Logger was: '%s'. Expected: 'msg=info'", msg)
	assert.Assert(t, strings.Contains(msg, "requestID=123"), "Logger was: '%s'. Expected: 'requestID=123'", msg)
}

func TestLogSomeRealMessages(t *testing.T) {

	// You should see these printed out, all correctly formatted.
	err := log.Debug("details", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	assert.Assert(t, err == nil)

	err = log.Print("info", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	assert.Assert(t, err == nil)

	err = log.Error(errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	assert.Assert(t, err == nil)

	scope := log.WithScope(log.Fields{"scopeID": 123})
	assert.Assert(t, scope != nil)

	err = scope.Debug("details", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	assert.Assert(t, err == nil)

	err = scope.Print("info", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	assert.Assert(t, err == nil)

	err = scope.Error(errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	assert.Assert(t, err == nil)
}

func BenchmarkLogging(b *testing.B) {
	logger := log.New(func(conf *log.Config) {
		conf.Output = ioutil.Discard
	})

	fields := log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	}

	for n := 0; n < b.N; n++ {
		_ = logger.Print("test details", fields)
	}

}