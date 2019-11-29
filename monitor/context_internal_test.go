package monitor

import (
	"context"
	"testing"

	"gotest.tools/assert"
)

func TestContext_Internal_Success(t *testing.T) {
	txn := &Transaction{}

	ctx := context.TODO()
	ctx = context.WithValue(ctx, txnContextKey, txn)
	txn, err := TxnFromContext(ctx)

	assert.Assert(t, txn != nil, txn)
	assert.Assert(t, err == nil, err)
}
