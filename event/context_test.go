package event_test

import (
	"context"
	"testing"

	"github.com/cultureamp/glamplify/event"
	"gotest.tools/assert"
)

func TestContext_Fail(t *testing.T) {

	ctx := context.TODO()
	txn, err := event.TxnFromContext(ctx)

	assert.Assert(t, txn == nil, txn)
	assert.Assert(t, err != nil, err)
}
