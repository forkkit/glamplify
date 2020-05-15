package context

import (
	"github.com/cultureamp/glamplify/aws"
	"github.com/cultureamp/glamplify/jwt"
	"net/http"
)

func GetRequestScopedFieldsFromRequest(r *http.Request) (RequestScopedFields, bool) {
	return GetRequestScopedFieldsFromCtx(r.Context())
}

func AddRequestScopedFieldsRequest(r *http.Request, requestScopeFields RequestScopedFields) *http.Request {
	ctx := AddRequestScopedFieldsToCtx(r.Context(), requestScopeFields)
	return r.WithContext(ctx)
}

// WrapRequest returns the same *http.Request if TraceID was already present in the context.
// If TraceID was missing, then it checks xray and if present, adds that TraceID, or if missing creates a new TraceID.
// If a TraceID was added (from xray or new) to the context, then this method also tries to decode the JWT payload and
// adds CustomerAggregateID and UserAggregateID if successful.
func WrapRequest(r *http.Request) *http.Request {
	rsFields, ok := GetRequestScopedFieldsFromRequest(r)
	if ok {
		return r
	}

	// need to create new RequestScopedFields
	ctx := r.Context()
	traceID, _ := aws.GetTraceID(ctx) // creates new TraceID if xray hasn't already added to the context
	requestID := traceID
	payload, err := jwt.PayloadFromRequest(r)

	if err == nil {
		rsFields = NewRequestScopeFields(traceID, requestID, payload.Customer, payload.EffectiveUser)
	} else {
		rsFields = NewRequestScopeFields(traceID, requestID, "", "")
	}
	ctx = rsFields.AddToCtx(ctx)
	return r.WithContext(ctx)
}

// WrapRequestWithDecoder returns the same *http.Request if TraceID was already present in the context.
// If TraceID was missing, then it checks xray and if present, adds that TraceID, or if missing creates a new TraceID.
// If a TraceID was added (from xray or new) to the context, then this method also tries to decode the JWT payload and
// adds CustomerAggregateID and UserAggregateID if successful.
func WrapRequestWithDecoder(r *http.Request, jwtDecoder jwt.DecodeJwtToken) *http.Request {
	rsFields, ok := GetRequestScopedFieldsFromRequest(r)
	if ok {
		return r
	}

	// need to create new RequestScopedFields
	ctx := r.Context()
	traceID, _ := aws.GetTraceID(ctx) // creates new TraceID if xray hasn't already added to the context
	requestID := traceID
	payload, err := jwt.PayloadFromRequestWithDecoder(r, jwtDecoder)

	if err == nil {
		rsFields = NewRequestScopeFields(traceID, requestID, payload.Customer, payload.EffectiveUser)
	} else {
		rsFields = NewRequestScopeFields(traceID, requestID,"", "")
	}
	ctx = rsFields.AddToCtx(ctx)
	return r.WithContext(ctx)
}