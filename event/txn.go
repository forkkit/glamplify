package event

import (
	"context"
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

// Transaction is a wrapper over the underlying implementation
type Transaction struct {
	impl    newrelic.Transaction
	logging bool
	logger  *eventLogger
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
			txn.logError("AddAtributes", err)
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
	count, err := txn.impl.Write(bytes)
	txn.logError("Write", err)
	return count, err
}

// WriteHeader delegates to the wrapped response
func (txn Transaction) WriteHeader(statusCode int) {
	txn.impl.WriteHeader(statusCode)
}

func (txn Transaction) logError(msg string, err error) {
	if err != nil && txn.logging {
		txn.logger.Error(msg, map[string]interface{}{
			"error": err,
		})
	}
}

// FromContext todo
func FromContext(ctx context.Context) *Transaction {
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		return &Transaction{
			impl:    txn,
			logging: false,
			logger:  nil,
		}
	}

	// TODO log error!
	return nil

}
