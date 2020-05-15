package context_test

import (
	"context"
	context2 "github.com/cultureamp/glamplify/context"
	"testing"

	"gotest.tools/assert"
)

func Test_TransactionFields_New(t *testing.T) {
	transactionFields := context2.NewRequestScopeFields("1-2-3", "7-8-9", "hooli", "UserAggregateID-123")
	assert.Assert(t, transactionFields.TraceID == "1-2-3",  transactionFields.TraceID)
	assert.Assert(t, transactionFields.RequestID == "7-8-9",  transactionFields.RequestID)
	assert.Assert(t, transactionFields.CustomerAggregateID == "hooli",  transactionFields.CustomerAggregateID)
	assert.Assert(t, transactionFields.UserAggregateID == "UserAggregateID-123",  transactionFields.UserAggregateID)
}

func Test_TransactionFields_NewFromCtx(t *testing.T) {
	transactionFields := context2.NewRequestScopeFields("1-2-3", "7-8-9","hooli", "UserAggregateID-123")

	ctx := context.Background()
	ctx = transactionFields.AddToCtx(ctx)

	rsFields, ok := context2.GetRequestScopedFieldsFromCtx(ctx)
	assert.Assert(t, ok, ok)
	assert.Assert(t, rsFields.TraceID == "1-2-3",  transactionFields.TraceID)
	assert.Assert(t, rsFields.RequestID == "7-8-9",  transactionFields.RequestID)
	assert.Assert(t, rsFields.CustomerAggregateID == "hooli",  transactionFields.CustomerAggregateID)
	assert.Assert(t, rsFields.UserAggregateID == "UserAggregateID-123",  transactionFields.UserAggregateID)
}

