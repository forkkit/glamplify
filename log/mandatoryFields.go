package log

import "context"

type MandatoryFields struct {
	TraceId             string `json:"trace_id"`
	UserAggregateId     string `json:"user"`
	CustomerAggregateId string `json:"customer"`
}

func NewMandatoryFields(traceId string, customer string, user string) MandatoryFields {
	return MandatoryFields{
		TraceId: traceId,
		CustomerAggregateId: customer,
		UserAggregateId: user,
	}
}

func NewMandatoryFieldsFromCtx(ctx context.Context) MandatoryFields {
	cfg := MandatoryFields{}

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

func (mFields MandatoryFields) AddToCtx(ctx context.Context) context.Context {
	ctx = AddTraceId(ctx, mFields.TraceId)
	ctx = AddCustomer(ctx, mFields.CustomerAggregateId)
	return AddUser(ctx, mFields.UserAggregateId)
}