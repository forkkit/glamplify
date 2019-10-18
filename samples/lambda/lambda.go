package main

import (
	"context"

	"github.com/cultureamp/glamplify/event"
	"github.com/cultureamp/glamplify/log"
)

func main() {
	app, _ := event.NewApplication("Glamplify-Lambda-Demo", func(conf *event.Config) {
		conf.Enabled = true        // default = "false"
		conf.Logging = true        // default = "false"
		conf.ServerlessMode = true // default = "false"
	})

	event.Start(handler, app)
}

func handler(ctx context.Context, input string) (string, error) {
	log.Print("Begin handler", log.Fields{
		"input": input,
	})

	app, err := event.AppFromContext(ctx)
	if err != nil {
		return "APP ERROR", err
	}

	txn, err := event.TxnFromContext(ctx)
	if err != nil {
		return "TXN ERROR", err
	}

	log.Print("End handler", log.Fields{
		"app": app,
		"txn": txn,
	})

	return "Ok", nil
}
