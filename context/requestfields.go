package context

import (
	"context"
)

type RequestScopedFields struct {
	TraceID             string `json:"trace_id"`
	RequestID           string `json:"request_id"`
	UserAggregateID     string `json:"user"`
	CustomerAggregateID string `json:"customer"`
}

func NewRequestScopeFields(traceID string, requestID string, customer string, user string) RequestScopedFields {
	return RequestScopedFields{
		TraceID:             traceID,
		RequestID:           requestID,
		CustomerAggregateID: customer,
		UserAggregateID:     user,
	}
}

func (rsFields RequestScopedFields) AddToCtx(ctx context.Context) context.Context {
	ctx = AddTraceID(ctx, rsFields.TraceID)
	ctx = AddRequestID(ctx, rsFields.RequestID)
	ctx = AddCustomer(ctx, rsFields.CustomerAggregateID)
	return AddUser(ctx, rsFields.UserAggregateID)
}
