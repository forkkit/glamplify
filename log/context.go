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

func ConfigToCtx(ctx context.Context, cfg Config) context.Context {
	ctx = AddTraceId(ctx, cfg.TraceId)
	ctx = AddCustomer(ctx, cfg.CustomerAggregateId)
	return AddUser(ctx, cfg.UserAggregateId)
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
func ConfigFromCtx(ctx context.Context) Config {
	cfg := Config{}

	if prod, ok := GetTraceId(ctx); ok {
		cfg.TraceId = prod
	}
	if prod, ok := GetCustomer(ctx); ok {
		cfg.CustomerAggregateId = prod
	}
	if prod, ok := GetUser(ctx); ok {
		cfg.UserAggregateId = prod
	}

	return cfg
}