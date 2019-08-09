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
	assertKV(t, msg, "msg=details")
	assertKV(t, msg, "severity=DEBUG")
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
	assertKV(t, msg, "msg=details")
	assertKV(t, msg, "severity=DEBUG")
	assertKV(t, msg, "string=hello")
	assertKV(t, msg, "int=123")
	assertKV(t, msg, "float=42.48")
	assertKV(t, msg, "string2=\"hello world\"")
	assertKV(t, msg, "\"string3 space\"=world")
}

func TestPrint_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	err := logger.Print("info")
	assert.Assert(t, err == nil)

	msg := memBuffer.String()
	assertKV(t, msg, "msg=info")
	assertKV(t, msg, "severity=INFO")
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
	assertKV(t, msg, "msg=info")
	assertKV(t, msg, "severity=INFO")
	assertKV(t, msg, "string=hello")
	assertKV(t, msg, "int=123")
	assertKV(t, msg, "float=42.48")
	assertKV(t, msg, "string2=\"hello world\"")
	assertKV(t, msg, "\"string3 space\"=world")
}

func TestError_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	err := logger.Error(errors.New("error"))
	assert.Assert(t, err == nil)

	msg := memBuffer.String()
	assertKV(t, msg, "error=error")
	assertKV(t, msg, "severity=ERROR")
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
	assertKV(t, msg, "error=error")
	assertKV(t, msg, "severity=ERROR")
	assertKV(t, msg, "string=hello")
	assertKV(t, msg, "int=123")
	assertKV(t, msg, "float=42.48")
	assertKV(t, msg, "string2=\"hello world\"")
	assertKV(t, msg, "\"string3 space\"=world")
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
	assertKV(t, msg, "msg=details")
	assertKV(t, msg, "requestID=123")

	memBuffer.Reset()

	err = scope.Print("info")
	assert.Assert(t, err == nil)

	msg = memBuffer.String()
	assertKV(t, msg, "msg=info")
	assertKV(t, msg, "requestID=123")
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

func assertKV(t *testing.T, log string, kv string) {
	assert.Assert(t, strings.Contains(log, kv), "Expected '%s' in '%s'", kv, log)
}
