package log_test

import (
	"context"
	"github.com/cultureamp/glamplify/log"
	"testing"

	"gotest.tools/assert"
)

func Test_TransactionFields_New(t *testing.T) {
	transactionFields := log.NewRequestScopeFields("1-2-3", "unilever", "UserAggregateID-123")
	assert.Assert(t, transactionFields.TraceID == "1-2-3",  transactionFields.TraceID)
	assert.Assert(t, transactionFields.CustomerAggregateID == "unilever",  transactionFields.CustomerAggregateID)
	assert.Assert(t, transactionFields.UserAggregateID == "UserAggregateID-123",  transactionFields.UserAggregateID)
}

func Test_TransactionFields_NewFromCtx(t *testing.T) {
	transactionFields := log.NewRequestScopeFields("1-2-3", "unilever", "UserAggregateID-123")

	ctx := context.Background()
	ctx = transactionFields.AddToCtx(ctx)

	transactionFields = log.NewRequestScopeFieldsFromCtx(ctx)
	assert.Assert(t, transactionFields.TraceID == "1-2-3",  transactionFields.TraceID)
	assert.Assert(t, transactionFields.CustomerAggregateID == "unilever",  transactionFields.CustomerAggregateID)
	assert.Assert(t, transactionFields.UserAggregateID == "UserAggregateID-123",  transactionFields.UserAggregateID)
}

