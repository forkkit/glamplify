package log

import (
	"context"
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

func (rsFields RequestScopedFields) AddToCtx(ctx context.Context) context.Context {
	ctx = AddTraceID(ctx, rsFields.TraceID)
	ctx = AddCustomer(ctx, rsFields.CustomerAggregateID)
	return AddUser(ctx, rsFields.UserAggregateID)
}


