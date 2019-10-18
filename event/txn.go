package event

import (
	"context"
	"errors"
	"net/http"

	"github.com/cultureamp/glamplify/log"
	newrelic "github.com/newrelic/go-agent"
)

// Transaction is a wrapper over the underlying implementation
type Transaction struct {
	impl    newrelic.Transaction
	app     *Application
	name    string
	logging bool
	logger  *eventLogger
}

func txnFromContext(ctx context.Context) (*Transaction, error) {
	txn, ok := ctx.Value(txnContextKey).(*Transaction)
	if ok && txn != nil {
		return txn, nil
	}

	err := errors.New("no transaction in context")
	return nil, err
}

// GetApplication todo
func (txn Transaction) GetApplication() *Application {
	return txn.app
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

	txn.log("Transaction Ended", log.Fields{
		"txnName": txn.name,
	})
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

func (txn *Transaction) addTransactionToHTTPContext(req *http.Request) *http.Request {
	ctx := txn.addTransactionToContext(req.Context())
	return req.WithContext(ctx)
}

func (txn *Transaction) addTransactionToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, txnContextKey, txn)
}

func (txn Transaction) logError(msg string, err error) {
	if err != nil && txn.logging {
		txn.logger.Error(msg, map[string]interface{}{
			"error": err,
		})
	}
}

func (txn Transaction) log(msg string, fields log.Fields) {
	if txn.logging {
		txn.logger.Debug(msg, fields)
	}
}
