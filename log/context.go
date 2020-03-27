package log

import (
	"context"
)

func AddTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, TraceIdCtx, traceId)
}

func AddCustomer(ctx context.Context, customer string) context.Context {
	return context.WithValue(ctx, CustomerCtx, customer)
}

func AddUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, UserCtx, user)
}