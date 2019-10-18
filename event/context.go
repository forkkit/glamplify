package event

import (
	"context"
	"errors"
	"net/http"
)

type key int

const (
	// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
	txnContextKey     key = iota
	appContextKey     key = iota
	handlerContextKey key = iota
)

// TxnFromRequest retrieves the current Transaction associated with the request, error is set appropriately
func TxnFromRequest(w http.ResponseWriter, r *http.Request) (*Transaction, error) {
	ctx := r.Context()
	return TxnFromContext(ctx)
}

// TxnFromContext retireves the current Transation from the given context, error is set appropriately
func TxnFromContext(ctx context.Context) (*Transaction, error) {
	txn, ok := ctx.Value(txnContextKey).(*Transaction)
	if ok && txn != nil {
		return txn, nil
	}

	err := errors.New("no transaction in context")
	return nil, err
}

// AppFromRequest retrives the current Application associated with the request, error is set appropriately
func AppFromRequest(w http.ResponseWriter, r *http.Request) (*Application, error) {
	ctx := r.Context()
	return AppFromContext(ctx)
}

// AppFromContext retireves the current Application from the given context, error is set appropriately
func AppFromContext(ctx context.Context) (*Application, error) {
	app, ok := ctx.Value(appContextKey).(*Application)
	if ok && app != nil {
		return app, nil
	}

	err := errors.New("no application in context")
	return nil, err
}

func handlerFromContext(ctx context.Context) (*lambdaHandler, error) {
	handler, ok := ctx.Value(handlerContextKey).(*lambdaHandler)
	if ok && handler != nil {
		return handler, nil
	}

	err := errors.New("no handler in context")
	return nil, err
}
