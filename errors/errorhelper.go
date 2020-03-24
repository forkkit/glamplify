package errors

import (
	"context"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/monitor"
	"github.com/cultureamp/glamplify/notify"
)

func HandleError(err error, fields log.Fields) {

	// call logger
	scope := log.FromScope(context.Background())
	scope.Error(err, log.Fields(fields))

	// call newrelic
	// todo - how to call new relic without a context?

	// call bugsnag
	notify.Error(err, fields)

}

func HandleErrorWithContext(ctx context.Context, err error, fields log.Fields) {

	// call logger
	scope := log.FromScope(ctx)
	scope.Error(err, fields)

	// call newrelic
	txn, err := monitor.TxnFromContext(ctx)
	if err == nil {
		txn.ReportErrorDetails(err.Error(), "app context", fields)
	}

	// call bugsnag
	notify.ErrorWithContext(err, ctx, fields)
}
