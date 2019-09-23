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

// AddAttributes adds customer data (key-value) to a current transaction (ie. http web request)
func (txn Transaction) AddAttributes(entries Entries) error {

	var err error
	for k, v := range entries {
		err = txn.impl.AddAttribute(k, v)
		if err != nil {
			return err
		}
	}
	return nil
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
