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

func ConfigFromCtx(ctx context.Context) Config {
	cfg := Config{}

	if prod, ok := ctx.Value(TraceIdCtx).(string); ok {
		cfg.TraceId = prod
	}

	if prod, ok := ctx.Value(CustomerCtx).(string); ok {
		cfg.Customer = prod
	}
	if prod, ok := ctx.Value(UserCtx).(string); ok {
		cfg.User = prod
	}

	return cfg
}