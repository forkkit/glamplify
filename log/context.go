package log

import (
	"context"
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
