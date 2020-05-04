package log

import (
	"context"
)

func AddTraceId(ctx context.Context, traceId string) context.Context {
	if len(traceId) > 0 {
		return context.WithValue(ctx, TraceIdCtx, traceId)
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

func GetTraceId(ctx context.Context) (string, bool) {
	traceId, ok := ctx.Value(TraceIdCtx).(string)
	return traceId, ok
}

func GetUser(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(UserCtx).(string)
	return user, ok
}

func GetCustomer(ctx context.Context) (string, bool) {
	customer, ok := ctx.Value(CustomerCtx).(string)
	return customer, ok
}
