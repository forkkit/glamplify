package log

import (
	"context"
	"github.com/cultureamp/glamplify/aws"
)

type RequestScopedFields struct {
	TraceID             string `json:"trace_id"`
	UserAggregateID     string `json:"user"`
	CustomerAggregateID string `json:"customer"`
}

func NewRequestScopeFields(traceID string, customer string, user string) RequestScopedFields {
	return RequestScopedFields{
		TraceID:             traceID,
		CustomerAggregateID: customer,
		UserAggregateID:     user,
	}
}

func NewRequestScopeFieldsFromCtx(ctx context.Context) RequestScopedFields {
	rsFields := RequestScopedFields{}

	// First look for our key in the context
	if prod, ok := GetTraceID(ctx); ok {
		rsFields.TraceID = prod
	} else {
		// if ours isn't there, then ask AWS SDK for the TraceID
		prod, ok = aws.GetTraceID(ctx)	// creates new traceID if missing, but doesn't add to ctx!
		rsFields.TraceID = prod
	}

	if prod, ok := GetCustomer(ctx); ok {
		rsFields.CustomerAggregateID = prod
	}
	if prod, ok := GetUser(ctx); ok {
		rsFields.UserAggregateID = prod
	}

	return rsFields
}

func (mFields RequestScopedFields) AddToCtx(ctx context.Context) context.Context {
	ctx = AddTraceID(ctx, mFields.TraceID)
	ctx = AddCustomer(ctx, mFields.CustomerAggregateID)
	return AddUser(ctx, mFields.UserAggregateID)
}