package monitor

import (
	"context"
	"os"
	"testing"
	"time"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
)

var (
	ctx      context.Context
	rsFields gcontext.RequestScopedFields
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	ctx = context.Background()
	ctx = gcontext.AddRequestFields(ctx, gcontext.RequestScopedFields{
		TraceID:             "1-2-3",
		RequestID:           "7-8-9",
		CorrelationID:       "1-5-9",
		CustomerAggregateID: "hooli",
		UserAggregateID:     "UserAggregateID-123",
	})

	rsFields, _ = gcontext.GetRequestScopedFields(ctx)

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


func Test_Realworld(t *testing.T) {
	// https://log-api.newrelic.com/log/v1
	writer := NewWriter(func(config *WriterConfig) {
		config.Endpoint ="https://log-api.newrelic.com/log/v1"
	})
	mlog := log.NewFromCtxWithCustomerWriter(ctx, writer)

	mlog.Info("hello_world2")
	time.Sleep(2 * time.Second)
}

