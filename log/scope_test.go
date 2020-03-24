package log

import (
	"bytes"
	"errors"
	"fmt"
	"gotest.tools/assert"
	"strings"
	"testing"
)


func TestScope(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx, Fields{
		"requestID": 123,
	})

	scope.Debug("detail_event")
	msg := memBuffer.String()
	assertScopeContainsString(t, msg, "event", "detail_event")
	assertScopeContainsInt(t, msg, "requestID", 123)

	memBuffer.Reset()
	scope.Info("info_event")
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "info_event")
	assertScopeContainsInt(t, msg, "requestID", 123)

	memBuffer.Reset()
	scope.Warn("warn_event")
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "warn_event")
	assertScopeContainsInt(t, msg, "requestID", 123)

	memBuffer.Reset()
	scope.Error(errors.New("error"))
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "error")
	assertScopeContainsInt(t, msg, "requestID", 123)

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "fatal")
			assertContainsString(t, msg, "severity", "FATAL")
		}
	}()

	scope.Fatal(errors.New("fatal")) // will call panic!
}

func TestScope_Overwrite(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx, Fields{
		"requestID": 123,
	})

	scope.Debug("detail_event", Fields {
		"requestID": 456,
	})
	msg := memBuffer.String()
	assertScopeContainsString(t, msg, "event", "detail_event")
	assertScopeContainsInt(t, msg, "requestID", 456)

	memBuffer.Reset()
	scope.Info("info_event", Fields {
		"requestID": 456,
	})
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "info_event")
	assertScopeContainsInt(t, msg, "requestID", 456)

	memBuffer.Reset()
	scope.Warn("warn_event", Fields {
		"requestID": 456,
	})
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "warn_event")
	assertScopeContainsInt(t, msg, "requestID", 456)

	memBuffer.Reset()
	scope.Error(errors.New("error"), Fields {
		"requestID": 456,
	})
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "error")
	assertScopeContainsInt(t, msg, "requestID", 456)

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertScopeContainsString(t, msg, "event", "fatal")
			assertScopeContainsString(t, msg, "severity", "FATAL")
			assertScopeContainsInt(t, msg, "requestID", 456)
		}
	}()

	// will call panic!
	scope.Fatal(errors.New("fatal"), Fields {
		"requestID": 456,
	})
}

func Test_RealWorld_Scope(t *testing.T) {

	scope := WithScope(ctx, Fields{"scopeID": 123})
	assert.Assert(t, scope != nil)

	scope.Debug("detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	scope.Info("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	scope.Warn("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	scope.Error(errors.New("error"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	defer func() {
		recover()
	}()

	// this will call panic!
	scope.Fatal(errors.New("fatal"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
}


func assertScopeContainsString(t *testing.T, log string, key string, val string) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":\"%s\"", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertScopeContainsInt(t *testing.T, log string, key string, val int) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":%v", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertScopeContainsSubDoc(t *testing.T, log string, key string, val string) {
	find := fmt.Sprintf("\"%s\":{\"%s\"", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

