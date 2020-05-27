package context

import (
	"context"
)

func AddRequestFields(ctx context.Context, rsFields RequestScopedFields) context.Context {
	return context.WithValue(ctx, RequestFieldsCtx, rsFields)
}

func GetRequestScopedFields(ctx context.Context) (RequestScopedFields, bool) {
	rsFields, ok := ctx.Value(RequestFieldsCtx).(RequestScopedFields)
	return rsFields, ok
}
