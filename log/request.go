package log

import "net/http"

func SeedRequestWithRequestScopeFields(r *http.Request) *http.Request {
	ctx := SeedCtxWithRequestScopeFields(r.Context())
	return r.WithContext(ctx)
}

func GetRequestScopedFieldsRequest(r *http.Request) RequestScopedFields {
	return GetRequestScopedFieldsCtx(r.Context())
}

func AddRequestScopedFieldsRequest(r *http.Request, requestScopeFields RequestScopedFields) *http.Request {
	ctx := AddRequestScopedFieldsCtx(r.Context(), requestScopeFields)
	return r.WithContext(ctx)
}
