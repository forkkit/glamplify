package context

import (
	"context"
)

type RequestScopedFields struct {
	TraceID             string `json:"trace_id"`
	RequestID           string `json:"request_id"`
	CorrelationID       string `json:"correlation_id"`
	UserAggregateID     string `json:"user"`
	CustomerAggregateID string `json:"customer"`
}

func NewRequestScopeFields(traceID string, requestID string, correlationID string, customer string, user string) RequestScopedFields {
	return RequestScopedFields{
		TraceID:             traceID,
		RequestID:           requestID,
		CorrelationID:       correlationID,
		CustomerAggregateID: customer,
		UserAggregateID:     user,
	}
}

func (rsFields RequestScopedFields) AddToCtx(ctx context.Context) context.Context {
	return AddRequestFields(ctx, rsFields)
}
