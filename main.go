package main

import (
	"context"
	"errors"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/cultureamp/glamplify/config"
	http2 "github.com/cultureamp/glamplify/http"
	"github.com/cultureamp/glamplify/jwt"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/monitor"
	"github.com/cultureamp/glamplify/notify"
	"net/http"
	"net/http/httptest"
)

func main() {

	/* CONFIG */

	// settings will contain configuration data as read in from the config file.
	settings := config.Load()

	// Or if you want to look for a config file from a specific location use
	//settings = config.LoadFrom([]string{"${HOME}/settings"}, "config")

	// Then you can use
	if settings.App.Version > 2.0 {
		// to do
	}

	/* LOGGING */
	// Creating loggers is cheap. Create them on every request/run
	// DO NOT CACHE/REUSE THEM
	cfg := log.MandatoryFields{
		TraceId:             "abc",   // Get TraceID from context or from wherever you have it stored
		UserAggregateId:     "user1", // Get UserAggregateId from context or from wherever you have it stored
		CustomerAggregateId: "cust1", // Get CustomerAggregateId from context or from wherever you have it stored
	}
	logger := log.New(cfg)

	// or if you want a field to be present on each subsequent logging call do this:
	logger = log.New(cfg, log.Fields{"request_id": 123})

	/* Monitor & Notify */
	app, appErr := monitor.NewApplication("GlamplifyUnitTests", func(conf *monitor.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
		conf.Labels = monitor.Labels{
			"asset":          log.Unknown,
			"classification": "restricted",
			"workload":       "development",
			"camp":           "amplify",
		}
	})
	if appErr != nil {
		logger.Fatal(appErr)
	}

	notifier, notifyErr := notify.NewNotifier("GlamplifyUnitTests", func (conf *notify.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	if notifyErr != nil {
		logger.Fatal(notifyErr)
	}

	pattern, handler := http2.WrapHTTPHandler(app, notifier, "/", rootRequestHandler)
	h := http.HandlerFunc(handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", pattern, nil)
	h.ServeHTTP(rr, req)

	app.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {

	// Do things

	/* REQUEST LOGGING */
	payload, err := jwt.PayloadFromRequest(r)

	// Create the logging config for this request
	cfg := log.MandatoryFields{
		TraceId:             xray.TraceID(r.Context()), // Get TraceID from context or from wherever you have it stored
		UserAggregateId:     payload.EffectiveUser, 	// Get UserAggregateId from context or from wherever you have it stored
		CustomerAggregateId: payload.Customer,      	// Get CustomerAggregateId from context or from wherever you have it stored
	}

	// Then create a logger that will use those config values when writing out logs
	logger := log.New(cfg)
	logger.Debug("Something happened")

	// or use the default logger with config
	log.Debug(cfg, "something happened")

	// Emit debug trace with types
	// Fields can contain any type of variables
	logger.Debug("Something else happened", log.Fields{
		"aString": "hello",
		"aInt":    123,
		"aFloat":  42.48,
	})

	// Emit normal logging (can add optional types if required)
	// Typically Print will be sent onto 3rd party aggregation tools (eg. Splunk)
	logger.Info("Executing main")

	// Emit Error (can add optional types if required)
	// Errors will always be sent onto 3rd party aggregation tools (eg. Splunk)
	err = errors.New("failed to save record to db")
	logger.Error(err)

	// Emit Fatal (can add optional types if required) and PANIC!
	// Fatal error will always be sent onto 3rd party aggregation tools (eg. Splunk)
	//err = errors.New("program died")
	//logger.Fatal(err)

	/* NEW RELIC TRANSACTION */
	txn, err := monitor.TxnFromRequest(w, r)
	if err == nil {
		txn.AddAttributes(log.Fields{
			"aString": "hello world",
			"aInt":    123,
		})
	}

	// Do more things

	/* NEW RELIC Add Attributes */
	if err == nil {
		txn.AddAttributes(log.Fields{
			"aString2": "goodbye",
			"aInt2":    456,
		})
	}
}


func doSomething(ctx context.Context) {
	cfg := log.MandatoryFields{
		TraceId:             xray.TraceID(ctx),        // Get TraceID from context or from wherever you have it stored
		UserAggregateId:     jwt.DecodeUser(jwtToken), // Get UserAggregateId from context or from wherever you have it stored
		CustomerAggregateId: "cust1",                  // Get CustomerAggregateId from context or from wherever you have it stored
	}

}