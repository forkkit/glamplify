package log_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestScope(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	scope := logger.WithScope(log.Fields{
		"requestID": 123,
	})

	ctx := context.Background()
	scope.Debug(ctx,"detail_event")

	msg := memBuffer.String()
	assertScopeContainsString(t, msg, "event", "detail_event")
	assertScopeContainsInt(t, msg, "requestID", 123)

	memBuffer.Reset()
	scope.Info(ctx,"info_event")

	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "info_event")
	assertScopeContainsInt(t, msg, "requestID", 123)

	memBuffer.Reset()
	scope.Warn(ctx,"warn_event")

	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "warn_event")
	assertScopeContainsInt(t, msg, "requestID", 123)

	memBuffer.Reset()
	scope.Error(ctx, errors.New("error"))

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
	scope.Fatal(ctx, errors.New("fatal")) // will call panic!
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

