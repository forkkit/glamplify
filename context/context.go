package context

import (
	"context"
	"github.com/cultureamp/glamplify/aws"
)

func AddRequestFields(ctx context.Context, rsFields RequestScopedFields) context.Context {
	return context.WithValue(ctx, RequestFieldsCtx, rsFields)
}

func GetRequestScopedFieldsFromCtx(ctx context.Context) (RequestScopedFields, bool) {
	rsFields, ok := ctx.Value(RequestFieldsCtx).(RequestScopedFields)
	return rsFields, ok
}

// WrapCtx returns modified context IF TraceID was added
// (returns the same context if TraceID was already present)
func WrapCtx(ctx context.Context) context.Context {
	rsFields, ok := GetRequestScopedFieldsFromCtx(ctx)
	if ok {
		return ctx
	}

	// need to create new RequestScopedFields
	traceID, _ := aws.GetTraceID(ctx) // creates new TraceID if xray hasn't already added to the context
	rsFields = RequestScopedFields{
		TraceID: traceID,
	}
	return AddRequestFields(ctx, rsFields)
}
