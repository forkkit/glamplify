package log

import (
	"context"
	"github.com/cultureamp/glamplify/constants"
)

func AddTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, constants.TraceIdCtx, traceId)
}

func AddCustomer(ctx context.Context, customer string) context.Context {
	return context.WithValue(ctx, constants.CustomerCtx, customer)
}

func AddUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, constants.UserCtx, user)
}