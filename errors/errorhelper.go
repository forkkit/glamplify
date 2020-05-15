package errors

import (
	"context"
	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/monitor"
	"github.com/cultureamp/glamplify/notify"
)

func HandleError(err error, fields log.Fields) {

	// call logger
	logger := log.New(gcontext.RequestScopedFields{})
	logger.Error("error", err, log.Fields(fields))

	// call newrelic
	// todo - how to call new relic without a context?

	// call bugsnag
	notify.Error(err, fields)
}

func HandleErrorWithContext(ctx context.Context, err error, fields log.Fields) {

	// call logger
	logger := log.NewFromCtx(ctx)
	logger.Error("error", err, fields)

	// call newrelic
	txn, err := monitor.TxnFromContext(ctx)
	if err == nil {
		txn.ReportErrorDetails(err.Error(), "app context", fields)
	}

	// call bugsnag
	notify.ErrorWithContext(ctx, err, fields)
}
