package log

import "context"

type TransactionFields struct {
	TraceID             string `json:"trace_id"`
	UserAggregateID     string `json:"user"`
	CustomerAggregateID string `json:"customer"`
}

func NewRequestScopeFields(traceID string, customer string, user string) TransactionFields {
	return TransactionFields{
		TraceID:             traceID,
		CustomerAggregateID: customer,
		UserAggregateID:     user,
	}
}

func NewRequestScopeFieldsFromCtx(ctx context.Context) TransactionFields {
	rsFields := TransactionFields{}

	if prod, ok := GetTraceID(ctx); ok {
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

func (mFields TransactionFields) AddToCtx(ctx context.Context) context.Context {
	ctx = AddTraceID(ctx, mFields.TraceID)
	ctx = AddCustomer(ctx, mFields.CustomerAggregateID)
	return AddUser(ctx, mFields.UserAggregateID)
}