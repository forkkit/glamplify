package log

import (
	"context"
	"github.com/cultureamp/glamplify/aws"
)

func AddTraceID(ctx context.Context, traceID string) context.Context {
	if len(traceID) > 0 {
		return context.WithValue(ctx, TraceIDCtx, traceID)
	}
	return ctx
}

func AddCustomer(ctx context.Context, customer string) context.Context {
	if len(customer) > 0 {
		return context.WithValue(ctx, CustomerCtx, customer)
	}
	return ctx
}

func AddUser(ctx context.Context, user string) context.Context {
	if len(user) > 0 {
		return context.WithValue(ctx, UserCtx, user)
	}
	return ctx
}

func GetTraceID(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(TraceIDCtx).(string)
	return traceID, ok
}

func GetUser(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(UserCtx).(string)
	return user, ok
}

func GetCustomer(ctx context.Context) (string, bool) {
	customer, ok := ctx.Value(CustomerCtx).(string)
	return customer, ok
}

func GetRequestScopedFieldsFromCtx(ctx context.Context) (RequestScopedFields, bool) {
	rsFields := RequestScopedFields{}

	val, ok := GetTraceID(ctx)
	if ok {
		rsFields.TraceID = val
	}
	val, ok = GetUser(ctx)
	if ok {
		rsFields.UserAggregateID = val
	}
	val, ok = GetCustomer(ctx)
	if ok {
		rsFields.CustomerAggregateID = val
	}

	// if there is a traceID, then we assume the context has RequestScopedFields
	if len(rsFields.TraceID) > 0 {
		return rsFields, true
	}
	return rsFields, false
}

// EnsureRequestScopedFieldsPresentInCtx returns modified context IF TraceID was added
// (returns the same context if TraceID was already present)
func EnsureRequestScopedFieldsPresentInCtx(ctx context.Context) context.Context {
	rsFields, ok := GetRequestScopedFieldsFromCtx(ctx)
	if ok {
		return ctx
	}

	// need to create new RequestScopedFields
	traceID, _ := aws.GetTraceID(ctx)	// creates new TraceID if xray hasn't already added to the context
	rsFields = NewRequestScopeFields(traceID,"","")
	return rsFields.AddToCtx(ctx)
}

func AddRequestScopedFieldsToCtx(ctx context.Context, requestScopeFields RequestScopedFields) context.Context {
	ctx = AddTraceID(ctx, requestScopeFields.TraceID)
	ctx = AddUser(ctx, requestScopeFields.UserAggregateID)
	return AddCustomer(ctx, requestScopeFields.CustomerAggregateID)
}