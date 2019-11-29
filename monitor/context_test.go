package monitor_test

import (
	"context"
	"github.com/cultureamp/glamplify/monitor"
	"testing"

	"gotest.tools/assert"
)

func TestContext_Fail(t *testing.T) {

	ctx := context.TODO()
	txn, err := monitor.TxnFromContext(ctx)

	assert.Assert(t, txn == nil, txn)
	assert.Assert(t, err != nil, err)
}
