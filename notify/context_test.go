package notify_test

import (
	"context"
	"github.com/cultureamp/glamplify/notify"
	"gotest.tools/assert"
	"testing"
)

func TestContext_Fail(t *testing.T) {

	ctx := context.TODO()
	txn, err := notify.NotifyFromContext(ctx)

	assert.Assert(t, txn == nil, txn)
	assert.Assert(t, err != nil, err)
}
