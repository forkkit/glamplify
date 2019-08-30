package event

import (
	"context"
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

// Transaction is a wrapper over the underlying implementation
type Transaction struct {
	impl newrelic.Transaction
}

func (txn Transaction) addTransactionContext(req *http.Request) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, txnContextKey, txn)
	return req.WithContext(ctx)
}

// AddAttribute adds customer data (key-value) to a current transaction (ie. http web request)
func (txn Transaction) AddAttribute(key string, value interface{}) error {
	return txn.impl.AddAttribute(key, value)
}

// End closes the current transaction
func (txn Transaction) End() {
	txn.impl.End()
}

// Header delegates to the wrapped response
func (txn Transaction) Header() http.Header {
	return txn.impl.Header()
}

// Write delegates to the wrapped response
func (txn Transaction) Write(bytes []byte) (int, error) {
	return txn.impl.Write(bytes)
}

// WriteHeader delegates to the wrapped response
func (txn Transaction) WriteHeader(statusCode int) {
	txn.impl.WriteHeader(statusCode)
}
