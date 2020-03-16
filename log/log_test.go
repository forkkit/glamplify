package log_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cultureamp/glamplify/constants"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
)

func TestDebug_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	ctx := context.Background()
	logger.Debug(ctx, "detail_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "detail_event")
	assertContainsString(t, msg, "severity", "DEBUG")
}

func TestDebugWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	ctx := context.Background()
	logger.Debug(ctx, "detail_event", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "detail_event")
	assertContainsString(t, msg, "severity", "DEBUG")
	assertContainsString(t, msg, "string", "hello")
	assertContainsInt(t, msg, "int", 123)
	assertContainsFloat(t, msg, "float", 42.48)
	assertContainsString(t, msg, "string2", "hello world")
	assertContainsString(t, msg, "string3 space", "world")
}

func TestInfo_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	ctx := context.Background()
	logger.Info(ctx, "info_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "info_event")
	assertContainsString(t, msg, "severity", "INFO")
}

func TestInfoWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	ctx := context.Background()
	logger.Info(ctx, "info_event", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "info_event")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "string", "hello")
	assertContainsInt(t, msg, "int", 123)
	assertContainsFloat(t, msg, "float", 42.48)
	assertContainsString(t, msg, "string2", "hello world")
	assertContainsString(t, msg, "string3 space", "world")
}

func TestInfoWithDuplicateFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	ctx := context.Background()
	logger.Info(ctx, "info_event", log.Fields{
		constants.ResourceLogField: "res_id", // set a standard types, this should overwrite the default
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "info_event")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "resource", "res_id") // by default this would normally be set to "host"
}

func TestWarn_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	ctx := context.Background()
	logger.Warn(ctx, "warn_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "warn_event")
	assertContainsString(t, msg, "severity", "WARN")
}

func TestWarnWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	ctx := context.Background()
	logger.Warn(ctx, "warn_event", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "warn_event")
	assertContainsString(t, msg, "severity", "WARN")
	assertContainsString(t, msg, "string", "hello")
	assertContainsInt(t, msg, "int", 123)
	assertContainsFloat(t, msg, "float", 42.48)
	assertContainsString(t, msg, "string2", "hello world")
	assertContainsString(t, msg, "string3 space", "world")
}

func TestError_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	ctx := context.Background()
	logger.Error(ctx, errors.New("error"))

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "error")
	assertContainsString(t, msg, "severity", "ERROR")
}

func TestErrorWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	ctx := context.Background()
	logger.Error(ctx, errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "error")
	assertContainsString(t, msg, "severity", "ERROR")
	assertContainsString(t, msg, "string", "hello")
	assertContainsInt(t, msg, "int", 123)
	assertContainsFloat(t, msg, "float", 42.48)
	assertContainsString(t, msg, "string2", "hello world")
	assertContainsString(t, msg, "string3 space", "world")
}

func TestFatal_Success(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "fatal")
			assertContainsString(t, msg, "severity", "FATAL")
		}
	}()

	ctx := context.Background()
	logger.Fatal(ctx, errors.New("fatal")) // will call panic!
}

func TestFatalWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "fatal")
			assertContainsString(t, msg, "severity", "FATAL")
			assertContainsString(t, msg, "string", "hello")
			assertContainsInt(t, msg, "int", 123)
			assertContainsFloat(t, msg, "float", 42.48)
			assertContainsString(t, msg, "string2", "hello world")
			assertContainsString(t, msg, "string3 space", "world")
		}
	}()

	// this will call panic!
	ctx := context.Background()
	logger.Fatal(ctx, errors.New("fatal"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
}

func TestNamespace_Success(t *testing.T) {

	t1 := time.Now()
	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer
	})

	time.Sleep(123 * time.Millisecond)
	t2 := time.Now()
	d := t2.Sub(t1)

	ctx := context.Background()
	logger.Error(ctx, errors.New("error"), log.Fields{
		"string": "hello",
		"int":    123,
		"float":  42.48,
		"reports_shared": log.Fields{
			"report":   "report1",
			"user":     "userid",
			"duration": fmt.Sprintf("P%gS", d.Seconds()),
		},
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "report", "report1")
	assertContainsString(t, msg, "user", "userid")

	assert.Assert(t, strings.Contains(msg, "reports_shared\":{\"duration"), "Expected 'reports_shared' in '%s'", msg)
}

func TestLogSomeRealMessages(t *testing.T) {

	ctx := context.Background()

	// You should see these printed out, all correctly formatted.
	log.Debug(ctx,"detail_event", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	log.Info(ctx, "info_event", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	log.Warn(ctx, "info_event", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	log.Error(ctx, errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	// multiple fields collections
	log.Info(ctx, "info_event", log.Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, log.Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	scope := log.WithScope(log.Fields{"scopeID": 123})
	assert.Assert(t, scope != nil)

	scope.Debug(ctx, "detail_event", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	scope.Info(ctx, "info_event", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	scope.Warn(ctx,"info_event", log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	scope.Error(ctx, errors.New("error"), log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
}

func Test_DurationAsIso8601(t *testing.T) {

	d := time.Millisecond * 456
	s := log.DurationAsISO8601(d)
	assert.Assert(t, s == "P0.456S", "was: %s", s)

	d = time.Millisecond * 1456
	s = log.DurationAsISO8601(d)
	assert.Assert(t, s == "P1.456S", "was: %s", s)
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

	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		logger.Info(ctx,"test details", fields)
	}
}

func assertContainsString(t *testing.T, log string, key string, val string) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":\"%s\"", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertContainsInt(t *testing.T, log string, key string, val int) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":%v", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertContainsFloat(t *testing.T, log string, key string, val float32) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":%v", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertContainsSubDoc(t *testing.T, log string, key string, val string) {
	find := fmt.Sprintf("\"%s\":{\"%s\"", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)

}
