package log_test

import (
	"context"
	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
	"testing"
)

func Test_Context_AddGet(t *testing.T) {

	ctx := context.Background()

	ctx = log.AddTraceID(ctx, "trace1")
	ctx = log.AddCustomer(ctx, "cust1")
	ctx = log.AddUser(ctx, "user1")

	id, ok := log.GetTraceID(ctx)
	assert.Assert(t, ok && id == "trace1", id)
	id, ok = log.GetCustomer(ctx)
	assert.Assert(t, ok && id == "cust1", id)
	id, ok = log.GetUser(ctx)
	assert.Assert(t, ok && id == "user1", id)
}

func Test_Context_TraceID_AddGet_Empty(t *testing.T) {

	ctx := context.Background()

	ctx = log.AddTraceID(ctx, "")
	ctx = log.AddCustomer(ctx, "")
	ctx = log.AddUser(ctx, "")

	id, ok := log.GetTraceID(ctx)
	assert.Assert(t, !ok && id == "", id)
	id, ok = log.GetTraceID(ctx)
	assert.Assert(t, !ok && id == "", id)
	id, ok = log.GetTraceID(ctx)
	assert.Assert(t, !ok && id == "", id)
}

func Test_Context_RequestScope_AddGet(t *testing.T) {

	rsFields := log.RequestScopedFields{
		TraceID: "1-2-3",
		UserAggregateID: "a-b-c",
		CustomerAggregateID: "xyz",
	}

	ctx := context.Background()
	ctx = log.AddRequestScopedFieldsToCtx(ctx, rsFields)

	resultFields, ok := log.GetRequestScopedFieldsFromCtx(ctx)
	assert.Assert(t, ok, ok)
	assert.Assert(t, resultFields.TraceID == rsFields.TraceID, resultFields.TraceID)
	assert.Assert(t, resultFields.CustomerAggregateID == rsFields.CustomerAggregateID, resultFields.CustomerAggregateID)
	assert.Assert(t, resultFields.UserAggregateID == rsFields.UserAggregateID, resultFields.UserAggregateID)
}


func Test_Context_Ensure(t *testing.T) {
	ctx := context.Background()

	ctx = log.EnsureRequestScopedFieldsPresentInCtx(ctx)
	id, ok := log.GetTraceID(ctx)
	assert.Assert(t, ok && id != "", id)
}
