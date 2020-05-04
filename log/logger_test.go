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

var (
	ctx               context.Context
	transactionFields TransactionFields
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	ctx = context.Background()
	ctx = AddTraceID(ctx, "1-2-3")
	ctx = AddCustomer(ctx, "unilever")
	ctx = AddUser(ctx, "UserAggregateID-123")

	transactionFields = NewRequestScopeFieldsFromCtx(ctx)

	os.Setenv("PRODUCT", "engagement")
	os.Setenv("APP", "murmur")
	os.Setenv("APP_VERSION", "87.23.11")
	os.Setenv("AWS_REGION", "us-west-02")
	os.Setenv("AWS_ACCOUNT_ID", "aws-account-123")
}

func shutdown() {
	os.Unsetenv("PRODUCT")
	os.Unsetenv("APP")
	os.Unsetenv("APP_VERSION")
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_ACCOUNT_ID")
}

func Test_New(t *testing.T) {
	logger := New(transactionFields)
	assert.Assert(t, logger != nil, logger)
}

func Test_NewWithContext(t *testing.T) {
	tempCtx, logger := NewWithCtx(ctx)
	assert.Assert(t, logger != nil, logger)

	traceID1, ok1 := GetTraceID(ctx)
	traceID2, ok2 := GetTraceID(tempCtx)

	assert.Assert(t, ok1, ok1)
	assert.Assert(t, ok2, ok2)
	assert.Assert(t, traceID1 == traceID2, "different trace ids")
	assert.Assert(t, traceID2 == "1-2-3", traceID2)
}

func Test_NewWithRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "*", nil)

	req1 := req.WithContext(ctx)
	req2, logger := NewWithRequest(req1)
	assert.Assert(t, logger != nil, logger)

	traceID1, ok1 := GetTraceID(req1.Context())
	traceID2, ok2 := GetTraceID(req2.Context())

	assert.Assert(t, ok1, ok1)
	assert.Assert(t, ok2, ok2)
	assert.Assert(t, traceID1 == traceID2, "different trace ids")
	assert.Assert(t, traceID2 == "1-2-3", traceID2)
}

func Test_Log_Debug(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

	logger.Debug( "detail_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "detail_event")
	assertContainsString(t, msg, "severity", "DEBUG")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_DebugWithFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

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
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_Info(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

	logger.Info("info_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "info_event")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_InfoWithFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

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
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_InfoWithDuplicateFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

	logger.Info("info_event", Fields{
		Resource: "res_id", // set a standard types, this should overwrite the default
	})

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "info_event")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "resource", "res_id") // by default this would normally be set to "host"
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_Warn(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

	logger.Warn("warn_event")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "warn_event")
	assertContainsString(t, msg, "severity", "WARN")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_WarnWithFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

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
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_Error(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

	logger.Error(errors.New("error"))

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "error")
	assertContainsString(t, msg, "severity", "ERROR")
	assertContainsString(t, msg, "trace_id", "1-2-3")
	assertContainsString(t, msg, "customer", "unilever")
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_ErrorWithFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

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
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")
}

func Test_Log_Fatal(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "fatal")
			assertContainsString(t, msg, "severity", "FATAL")
			assertContainsString(t, msg, "trace_id", "1-2-3")
			assertContainsString(t, msg, "customer", "unilever")
			assertContainsString(t, msg, "user", "UserAggregateID-123")
			assertContainsString(t, msg, "product", "engagement")
			assertContainsString(t, msg, "app", "murmur")
			assertContainsString(t, msg, "app_version", "87.23.11")
			assertContainsString(t, msg, "aws_region", "us-west-02")
			assertContainsString(t, msg, "aws_account_id", "aws-account-123")
		}
	}()

	logger.Fatal(errors.New("fatal")) // will call panic!
}

func Test_Log_FatalWithFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

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
			assertContainsString(t, msg, "user", "UserAggregateID-123")
			assertContainsString(t, msg, "product", "engagement")
			assertContainsString(t, msg, "app", "murmur")
			assertContainsString(t, msg, "app_version", "87.23.11")
			assertContainsString(t, msg, "aws_region", "us-west-02")
			assertContainsString(t, msg, "aws_account_id", "aws-account-123")
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

func Test_Log_Namespace(t *testing.T) {

	t1 := time.Now()
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer)

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
	assertContainsString(t, msg, "user", "UserAggregateID-123")
	assertContainsString(t, msg, "product", "engagement")
	assertContainsString(t, msg, "app", "murmur")
	assertContainsString(t, msg, "app_version", "87.23.11")
	assertContainsString(t, msg, "aws_region", "us-west-02")
	assertContainsString(t, msg, "aws_account_id", "aws-account-123")

	assertContainsSubDoc(t, msg, "reports_shared", "duration")
}


func TestScope(t *testing.T) {
	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer, Fields{
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
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(transactionFields, writer, Fields{
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

func Test_RealWorld(t *testing.T) {
	logger := New(transactionFields)

	// You should see these printed out, all correctly formatted.
	logger.Debug("detail_event", Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
	Debug(transactionFields, "detail_event", Fields{
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
	Info(transactionFields, "info_event", Fields{
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
	Warn(transactionFields, "info_event", Fields{
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
	Error(transactionFields, errors.New("error"), Fields{
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

	defer func() {
		recover()
	}()

	// this will call panic!
	Fatal(transactionFields, errors.New("fatal"), Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	})
}

func Test_RealWorld_Combined(t *testing.T) {
	logger := New(transactionFields)

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
	Debug(transactionFields, "detail_event", Fields{
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
	Info(transactionFields, "info_event", Fields{
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
	Warn(transactionFields, "warn_event", Fields{
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
	Error(transactionFields, errors.New("error"), Fields{
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

	defer func() {
		recover()
	}()

	// this will call panic!
	Fatal(transactionFields, errors.New("fatal"), Fields{
		"string1": "hello",
		"int1":    123,
		"float1":  42.48,
	}, Fields{
		"string2": "world",
		"int2":    456,
		"float2":  78.98,
	})
}



func Test_RealWorld_Scope(t *testing.T) {

	logger := New(transactionFields, Fields{"scopeID": 123})
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
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = ioutil.Discard
	})
	logger := newLogger(transactionFields, writer)

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