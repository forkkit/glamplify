package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gotest.tools/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var ctx context.Context

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

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

func Test_New(t *testing.T) {

	_, logger := New(ctx)
	assert.Assert(t, logger != nil, logger)

	req, _ := http.NewRequest("GET", "/", nil)
	_, logger, err := NewFromRequest(ctx, req)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, logger != nil, logger)

	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0In0.oDXzd1tq5XRpENEcw7GxAglOpRWmL5ld7XYPeNlrF-IfWYYRy86rta9yG9ug5wS1GV7Lvv8EbufXk0DKTnd23oObWoJtXLUaHh2TG9sw9bsxNwLKu1eWw7MQtUYByN2QFpRGeMQo_yw5Y6bT76janQ1NZknopHHvttcLBFuSMdThMX-4gOlaCuVsr8MQ218WUC-rVrSAol57at_2gf8PkEik3bcOd4bvUpf-ThumkljyzSrxVBY57H1kYbYAST4CwcCrf2F3oTLa_xNbFycngVCvJLZtSQR5GxwpO_ERqFziEaQ07bW6Svcs0EvARjCB-4vYdKTFaw3J5qu2aWVHf9m3a4QPA5O91ODYFYq_7k6upmxQl074_MQ-ZsnDRt0cUyPJjObMjU99MuMLQNnAMU67iNYkOxocR1OCNzLL1ObpeoYVq8sZWQPVhrPFDnC-V5uIsoSl9NofwcApLfUV2WjcMHxPfJYqPo-BNq3P_p1G1WSJ7iLP1BMXAU_ZaK49YaWb3fwu4NzRSCjsulWjMiE1yQL_bQrj4crygAyCgG7hpgq9OdiVl7YElrOL-oY1_3XCvnVcZkCd5dQjSbTXx-cW8Xc_zeY1QGxiKaeI3Yg24XLSVSFMNX4XNXwtNlK-LSrWQU8S0bVZBRDNo0jM9hx7INjYc4tamu2sGcH-71Q"
	authHeader := "Bearer " + token

	req.Header.Set("Authorization", authHeader)
	_, logger, err = NewFromRequest(ctx, req)
	// TODO - we get an error, because the "AUTH_PUBLIC_KEY" env var is not set to the public key
	// Need to inject this in somehow...
	//assert.Assert(t, err == nil, err)
	assert.Assert(t, logger != nil, logger)
}

func TestDebug_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	logger.Debug( "detail_event")

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
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	logger.Debug("detail_event", Fields{
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
	assertContainsString(t, msg, "string3_space", "world")
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
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	logger.Info("info_event")

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
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	logger.Info("info_event", Fields{
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
	assertContainsString(t, msg, "string3_space", "world")
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
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	logger.Info("info_event", Fields{
		Resource: "res_id", // set a standard types, this should overwrite the default
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
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	logger.Warn("warn_event")

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
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	logger.Warn("warn_event", Fields{
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
	assertContainsString(t, msg, "string3_space", "world")
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
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	logger.Error(errors.New("error"))

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
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	logger.Error(errors.New("error"), Fields{
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
	assertContainsString(t, msg, "string3_space", "world")
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
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

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

	logger.Fatal(errors.New("fatal")) // will call panic!
}

func TestFatalWithFields_Success(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "fatal")
			assertContainsString(t, msg, "severity", "FATAL")
			assertContainsString(t, msg, "string", "hello")
			assertContainsInt(t, msg, "int", 123)
			assertContainsFloat(t, msg, "float", 42.48)
			assertContainsString(t, msg, "string2", "hello world")
			assertContainsString(t, msg, "string3_space", "world")
			assertContainsString(t, msg, "trace_id", "1-2-3")
			assertContainsString(t, msg, "customer", "unilever")
			assertContainsString(t, msg, "user", "user-123")
			assertContainsString(t, msg, "product", "engagement")
			assertContainsString(t, msg, "app", "murmur")
			assertContainsString(t, msg, "app_version", "87.23.11")
			assertContainsString(t, msg, "region", "us-west-02")
		}
	}()

	logger.Fatal(errors.New("fatal"), Fields{ // this will call panic!
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
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer)

	time.Sleep(123 * time.Millisecond)
	t2 := time.Now()
	d := t2.Sub(t1)

	logger.Error(errors.New("error"), Fields{
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

	ctx = context.Background()
	_, logger := New(ctx)

	// You should see these printed out, all correctly formatted.
	logger.Debug("detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Info("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Warn("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Error(errors.New("error"), Fields{
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
	logger.Fatal(errors.New("fatal"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
}

func Test_RealWorld_Combined(t *testing.T) {

	ctx = context.Background()
	_, logger := New(ctx)

	// multiple fields collections
	logger.Debug("detail_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	logger.Info("info_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	logger.Warn("warn_event", Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})

	logger.Error(errors.New("error"), Fields{
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
	logger.Fatal(errors.New("fatal"), Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})
}


func TestScope(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer, Fields{
		"requestID": 123,
	})

	logger.Debug("detail_event")
	msg := memBuffer.String()
	assertScopeContainsString(t, msg, "event", "detail_event")
	assertScopeContainsInt(t, msg, "request_id", 123)

	memBuffer.Reset()
	logger.Info("info_event")
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "info_event")
	assertScopeContainsInt(t, msg, "request_id", 123)

	memBuffer.Reset()
	logger.Warn("warn_event")
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "warn_event")
	assertScopeContainsInt(t, msg, "request_id", 123)

	memBuffer.Reset()
	logger.Error(errors.New("error"))
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "error")
	assertScopeContainsInt(t, msg, "request_id", 123)

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "fatal")
			assertContainsString(t, msg, "severity", "FATAL")
		}
	}()

	logger.Fatal(errors.New("fatal")) // will call panic!
}

func TestScope_Overwrite(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *Config) {
		conf.Output = memBuffer
	})
	_, logger := NewWitCustomWriter(ctx, writer, Fields{
		"requestID": 123,
	})

	logger.Debug("detail_event", Fields {
		"requestID": 456,
	})
	msg := memBuffer.String()
	assertScopeContainsString(t, msg, "event", "detail_event")
	assertScopeContainsInt(t, msg, "request_id", 456)

	memBuffer.Reset()
	logger.Info("info_event", Fields {
		"requestID": 456,
	})
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "info_event")
	assertScopeContainsInt(t, msg, "request_id", 456)

	memBuffer.Reset()
	logger.Warn("warn_event", Fields {
		"requestID": 456,
	})
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "warn_event")
	assertScopeContainsInt(t, msg, "request_id", 456)

	memBuffer.Reset()
	logger.Error(errors.New("error"), Fields {
		"requestID": 456,
	})
	msg = memBuffer.String()
	assertScopeContainsString(t, msg, "event", "error")
	assertScopeContainsInt(t, msg, "request_id", 456)

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertScopeContainsString(t, msg, "event", "fatal")
			assertScopeContainsString(t, msg, "severity", "FATAL")
			assertScopeContainsInt(t, msg, "request_id", 456)
		}
	}()

	// will call panic!
	logger.Fatal(errors.New("fatal"), Fields {
		"request_id": 456,
	})
}

func Test_RealWorld_Scope(t *testing.T) {

	ctx = context.Background()
	_, logger := New(ctx, Fields{"scopeID": 123})
	assert.Assert(t, logger != nil)

	logger.Debug("detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Info("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Warn("info_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})

	logger.Error(errors.New("error"), Fields{
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
	logger.Fatal(errors.New("fatal"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
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
	writer := NewWriter(func(conf *Config) {
		conf.Output = ioutil.Discard
	})
	_, logger := newLogger(ctx, writer)

	fields := Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	}

	for n := 0; n < b.N; n++ {
		logger.Info("test details", fields)
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