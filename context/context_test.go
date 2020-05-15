package context_test

import (
	"context"
	context2 "github.com/cultureamp/glamplify/context"
	"gotest.tools/assert"
	"testing"
)

func Test_Context_AddGet(t *testing.T) {

	ctx := context.Background()

	ctx = context2.AddTraceID(ctx, "trace1")
	ctx = context2.AddRequestID(ctx, "request1")
	ctx = context2.AddCustomer(ctx, "cust1")
	ctx = context2.AddUser(ctx, "user1")

	id, ok := context2.GetTraceID(ctx)
	assert.Assert(t, ok && id == "trace1", id)
	id, ok = context2.GetRequestID(ctx)
	assert.Assert(t, ok && id == "request1", id)
	id, ok = context2.GetCustomer(ctx)
	assert.Assert(t, ok && id == "cust1", id)
	id, ok = context2.GetUser(ctx)
	assert.Assert(t, ok && id == "user1", id)
}

func Test_Context_TraceID_AddGet_Empty(t *testing.T) {

	ctx := context.Background()

	ctx = context2.AddTraceID(ctx, "")
	ctx = context2.AddRequestID(ctx, "")
	ctx = context2.AddCustomer(ctx, "")
	ctx = context2.AddUser(ctx, "")

	id, ok := context2.GetTraceID(ctx)
	assert.Assert(t, !ok && id == "", id)
	id, ok = context2.GetRequestID(ctx)
	assert.Assert(t, !ok && id == "", id)
	id, ok = context2.GetTraceID(ctx)
	assert.Assert(t, !ok && id == "", id)
	id, ok = context2.GetTraceID(ctx)
	assert.Assert(t, !ok && id == "", id)
}

func Test_Context_RequestScope_AddGet(t *testing.T) {

	rsFields := context2.RequestScopedFields{
		TraceID: "1-2-3",
		RequestID: "7-8-9",
		UserAggregateID: "a-b-c",
		CustomerAggregateID: "xyz",
	}

	ctx := context.Background()
	ctx = context2.AddRequestScopedFieldsToCtx(ctx, rsFields)

	resultFields, ok := context2.GetRequestScopedFieldsFromCtx(ctx)
	assert.Assert(t, ok, ok)
	assert.Assert(t, resultFields.TraceID == rsFields.TraceID, resultFields.TraceID)
	assert.Assert(t, resultFields.RequestID == rsFields.RequestID, resultFields.RequestID)
	assert.Assert(t, resultFields.CustomerAggregateID == rsFields.CustomerAggregateID, resultFields.CustomerAggregateID)
	assert.Assert(t, resultFields.UserAggregateID == rsFields.UserAggregateID, resultFields.UserAggregateID)
}


func Test_Context_Wrap(t *testing.T) {
	ctx := context.Background()

	ctx = context2.WrapCtx(ctx)
	id, ok := context2.GetTraceID(ctx)
	assert.Assert(t, ok && id != "", id)
}
