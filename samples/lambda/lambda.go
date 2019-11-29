package main

import (
	"context"
	"github.com/cultureamp/glamplify/helper"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/monitor"
)

func main() {
	app, _ := monitor.NewApplication("Glamplify-Lambda-Demo", func(conf *monitor.Config) {
		conf.Enabled = true        // default = "false"
		conf.Logging = true        // default = "false"
		conf.ServerlessMode = true // default = "false"
	})
	monitor.Start(handler, app)
}

func handler(ctx context.Context, input string) (string, error) {
	fields := log.Fields{
		"input": input,
	}
	log.Debug("Begin handler", fields)

	app, err := monitor.AppFromContext(ctx)
	if err != nil {
		helper.HandleError(err, fields)
		return "APP ERROR", err
	}

	txn, err := monitor.TxnFromContext(ctx)
	if err != nil {
		helper.HandleErrorWithContext(err, ctx, fields)
		return "TXN ERROR", err
	}

	log.Debug("End handler", log.Fields{
		"app": app,
		"txn": txn,
	})

	return "Ok", nil
}
