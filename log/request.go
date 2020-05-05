package log

import "net/http"

func SeedRequestCtxWithRequestScopeFields(r *http.Request) *http.Request {
	ctx := SeedCtxWithRequestScopeFields(r.Context())
	return r.WithContext(ctx)
}
