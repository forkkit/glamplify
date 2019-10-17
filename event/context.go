package event

import (
	"net/http"
)

type key int

const (
	// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
	txnContextKey key = iota
	appContextKey key = iota
)

// TxnFromRequest retrieves the current transaction associated with the request, ok is set appropriately
func TxnFromRequest(w http.ResponseWriter, r *http.Request) (Transaction, bool) {
	ctx := r.Context()
	txn, ok := ctx.Value(txnContextKey).(Transaction)
	return txn, ok
}

// AppFromRequest todo
func AppFromRequest(w http.ResponseWriter, r *http.Request) (Application, bool) {
	ctx := r.Context()
	app, ok := ctx.Value(appContextKey).(Application)
	return app, ok
}
