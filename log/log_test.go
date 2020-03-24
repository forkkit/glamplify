package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cultureamp/glamplify/constants"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

var ctx context.Context

func setup() {
	ctx = context.Background()
	ctx = AddTraceId(ctx, "1-2-3")
	ctx = AddCustomer(ctx, "unilever")
	ctx = AddUser(ctx, "user-123")

	os.Setenv("PRODUCT", "engagement")
	os.Setenv("APP", "murmur")
	os.Setenv("APP_VERSION", "87.23.11")
	os.Setenv("REGION", "us-west-02")
}

func shutdown() {
	os.Unsetenv("PRODUCT")
	os.Unsetenv("APP")
	os.Unsetenv("APP_VERSION")
	os.Unsetenv("REGION")
}

func TestDebug_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx)
	scope.Debug( "detail_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "detail_event")
	assertContainsString(t, msg, "severity", "DEBUG")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "user-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "region", "us-west-02")
}

func TestDebugWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx)
	scope.Debug("detail_event", Fields{
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
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "user-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "region", "us-west-02")
}

func TestInfo_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx)
	scope.Info("info_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "info_event")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "user-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "region", "us-west-02")
}

func TestInfoWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx)
	scope.Info("info_event", Fields{
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
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "user-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "region", "us-west-02")
}

func TestInfoWithDuplicateFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx)
	scope.Info("info_event", Fields{
		constants.ResourceLogField: "res_id", // set a standard types, this should overwrite the default
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "info_event")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "resource", "res_id") // by default this would normally be set to "host"
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "user-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "region", "us-west-02")
}

func TestWarn_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx)
	scope.Warn("warn_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "warn_event")
	assertContainsString(t, msg, "severity", "WARN")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "user-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "region", "us-west-02")
}

func TestWarnWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx)
	scope.Warn("warn_event", Fields{
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
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "user-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "region", "us-west-02")
}

func TestError_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx)
	scope.Error(errors.New("error"))

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "error")
	assertContainsString(t, msg, "severity", "ERROR")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "user-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "region", "us-west-02")
}

func TestErrorWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	scope := logger.withScope(ctx)
	scope.Error(errors.New("error"), Fields{
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
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "user-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "region", "us-west-02")
}

func TestFatal_Success(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "fatal")
			assertContainsString(t, msg, "severity", "FATAL")
			assertContainsString(t, msg, "trace_id", "1-2-3")
			assertContainsString(t, msg, "customer", "unilever")
			assertContainsString(t, msg, "user", "user-123")
			assertContainsString(t, msg, "product", "engagement")
			assertContainsString(t, msg, "app", "murmur")
			assertContainsString(t, msg, "app_version", "87.23.11")
			assertContainsString(t, msg, "region", "us-west-02")
		}
	}()

	scope := logger.withScope(ctx)
	scope.Fatal(errors.New("fatal")) // will call panic!
}

func TestFatalWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	logger := newLogger(func(conf *Config) {
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
			assertContainsString(t, msg, "trace_id", "1-2-3")
			assertContainsString(t, msg, "customer", "unilever")
			assertContainsString(t, msg, "user", "user-123")
			assertContainsString(t, msg, "product", "engagement")
			assertContainsString(t, msg, "app", "murmur")
			assertContainsString(t, msg, "app_version", "87.23.11")
			assertContainsString(t, msg, "region", "us-west-02")
		}
	}()

	scope := logger.withScope(ctx)
	scope.Fatal(errors.New("fatal"), Fields{  // this will call panic!
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
	logger := newLogger(func(conf *Config) {
		conf.Output = memBuffer
	})

	time.Sleep(123 * time.Millisecond)
	t2 := time.Now()
	d := t2.Sub(t1)

	scope := logger.withScope(ctx)
	scope.Error(errors.New("error"), Fields{
		"string": "hello",
		"int":    123,
		"float":  42.48,
		"reports_shared": Fields{
			"report":   "report1",
			"user":     "userid",
			"duration": fmt.Sprintf("P%gS", d.Seconds()),
		},
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "report", "report1")
	assertContainsString(t, msg, "user", "userid")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "user-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "region", "us-west-02")

	assertContainsSubDoc(t, msg, "reports_shared", "duration")
}

func Test_RealWorld(t *testing.T) {

	scope := WithScope(ctx)

	// You should see these printed out, all correctly formatted.
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

func Test_RealWorld_Combined(t *testing.T) {

	scope := WithScope(ctx)

	// multiple fields collections
	scope.Debug("detail_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	scope.Info("info_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	scope.Warn("warn_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	scope.Error(errors.New("error"), Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	defer func() {
		recover()
	}()

	// this will call panic!
	scope.Fatal(errors.New("fatal"), Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})
}

func Test_DurationAsIso8601(t *testing.T) {

	d := time.Millisecond * 456
	s := DurationAsISO8601(d)
	assert.Assert(t, s == "P0.456S", "was: %s", s)

	d = time.Millisecond * 1456
	s = DurationAsISO8601(d)
	assert.Assert(t, s == "P1.456S", "was: %s", s)
}

func BenchmarkLogging(b *testing.B) {
	logger := newLogger(func(conf *Config) {
		conf.Output = ioutil.Discard
	})

	fields := Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	}

	scope := logger.withScope(ctx)
	for n := 0; n < b.N; n++ {
		scope.Info("test details", fields)
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
