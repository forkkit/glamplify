package errors

import (
	"context"
	"github.com/cultureamp/glamplify/types"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/monitor"
	"github.com/cultureamp/glamplify/notify"
)

func HandleError(err error, fields types.Fields) {

	// call logger
	log.Error(err, log.Fields(fields))

	// call newrelic
	// todo - how to call new relic without a context?

	// call bugsnag
	notify.Error(err, fields)

}

func HandleErrorWithContext(err error, ctx context.Context, fields types.Fields) {

	// call logger
	log.Error(err, log.Fields(fields))

	// call newrelic
	txn, err := monitor.TxnFromContext(ctx)
	if err == nil {
		txn.ReportErrorDetails(err.Error(), "app context", monitor.Fields(fields))
	}

	// call bugsnag
	notify.ErrorWithContext(err, ctx, fields)
}
