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

	logger.Debug("details")

	msg := memBuffer.String()
	assertKV(t, msg, "msg=details")
	assertKV(t, msg, "severity=DEBUG")
	assertKV(t, msg, "forward-log=none")
}

func TestDebugWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	logger.Debug("details", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertKV(t, msg, "msg=details")
	assertKV(t, msg, "severity=DEBUG")
	assertKV(t, msg, "forward-log=none")
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

	logger.Print("info")

	msg := memBuffer.String()
	assertKV(t, msg, "msg=info")
	assertKV(t, msg, "severity=INFO")
	assertKV(t, msg, "forward-log=splunk")
}

func TestPrintWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	logger.Print("info", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertKV(t, msg, "msg=info")
	assertKV(t, msg, "severity=INFO")
	assertKV(t, msg, "forward-log=splunk")
	assertKV(t, msg, "string=hello")
	assertKV(t, msg, "int=123")
	assertKV(t, msg, "float=42.48")
	assertKV(t, msg, "string2=\"hello world\"")
	assertKV(t, msg, "\"string3 space\"=world")
}

func TestPrintWithDuplicateFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	logger.Print("info", log.Fields{
		log.FORWARD: "sumo", // set a standard field, this should overwrite the default
	})

	msg := memBuffer.String()
	assertKV(t, msg, "msg=info")
	assertKV(t, msg, "severity=INFO")
	assertKV(t, msg, "forward-log=sumo") // by default this would normally be set to "splunk"
}

func TestError_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	logger.Error(errors.New("error"))

	msg := memBuffer.String()
	assertKV(t, msg, "error=error")
	assertKV(t, msg, "severity=ERROR")
	assertKV(t, msg, "forward-log=splunk")
}

func TestErrorWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	logger.Error(errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertKV(t, msg, "error=error")
	assertKV(t, msg, "severity=ERROR")
	assertKV(t, msg, "forward-log=splunk")
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

	scope.Debug("details")

	msg := memBuffer.String()
	assertKV(t, msg, "forward-log=none")
	assertKV(t, msg, "msg=details")
	assertKV(t, msg, "requestID=123")

	memBuffer.Reset()

	scope.Print("info")

	msg = memBuffer.String()
	assertKV(t, msg, "forward-log=splunk")
	assertKV(t, msg, "msg=info")
	assertKV(t, msg, "requestID=123")
}

func TestLogSomeRealMessages(t *testing.T) {

	// You should see these printed out, all correctly formatted.
	log.Debug("details", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	log.Print("info", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	log.Error(errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	scope := log.WithScope(log.Fields{"scopeID": 123})
	assert.Assert(t, scope != nil)

	scope.Debug("details", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	scope.Print("info", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	scope.Error(errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
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
		logger.Print("test details", fields)
	}

}

func assertKV(t *testing.T, log string, kv string) {
	assert.Assert(t, strings.Contains(log, kv), "Expected '%s' in '%s'", kv, log)
}
