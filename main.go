package main

import (
	"context"
	"errors"
	"github.com/cultureamp/glamplify/aws"
	http2 "github.com/cultureamp/glamplify/http"
	"github.com/cultureamp/glamplify/jwt"
	"net/http"
	"net/http/httptest"

	"github.com/cultureamp/glamplify/config"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/monitor"
	"github.com/cultureamp/glamplify/notify"
)

func main() {

	// If you aren't passed a context, then you need to create a new one and then you should add
	// all the mandatory values to it so logging can retrieve them automatically
	// Example:
	ctx := context.Background()

	/* CONFIG */

	// settings will contain configuration data as read in from the config file.
	settings := config.Load()

	// Or if you want to look for a config file from a specific location use
	//settings = config.LoadFrom([]string{"${HOME}/settings"}, "config")

	// Then you can use
	if settings.App.Version > 2.0 {
		// to do
	}

	/* PARAMETER STORE / JWT */
	ps := aws.NewParameterStore(ctx, "default")
	pubKey, err := ps.Get("common/AUTH_PUBLIC_KEY")

	jwt := jwt.NewJWTDecoderFromBytes(ctx, []byte(pubKey))
	payload, err := jwt.Decode("")

	/* LOGGING */

	// AWS X-ray trace_id normally passed via http headers (_X_AMZN_TRACE_ID) or by another method
	// if you need to create a new one because you are the "start" of a tree then DON'T PASS/SET ANYTHING
	// and the logging system will create it automatically for you
	traceId :=  "1-58406520-a006649127e371903a2de979" // otherwise get it from header, etc
	ctx = log.AddTraceId(ctx, traceId)

	// If this service deals with a particularly customer, then set that on the context as well
	//customer := "FNSNDCJDF343"
	ctx = log.AddCustomer(ctx, payload.Customer)

	// And finally if this service deals with a particular user, then set that on the context as well
	//user := "JFOSNDJF97S"
	ctx = log.AddUser(ctx, payload.EffectiveUser)

	// Example below shows usage with the package level logger (sensible default)
	logger := log.New(ctx)

	// If you want to set some types for a particular scope (eg. for a Web Request
	// have a requestID for every log message for that logger) then you can use pass
	// log.Fields{} when creating the logger
	logger = log.New(ctx, log.Fields{"requestID": 123})

	// Emit debug trace
	// All messages must be static strings (as per Culture Amp Sensibile Default)
	logger.Debug("Something happened")

	// Emit debug trace with types
	// Fields can contain any type of variables
	logger.Debug("Something happened", log.Fields{
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
		logger.Error(appErr)
	}

	notifier, notifyErr := notify.NewNotifier("GlamplifyUnitTests", func (conf *notify.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	if notifyErr != nil {
		logger.Error(notifyErr)
	}

	pattern, handler := http2.WrapHTTPHandler(app, notifier, "/", rootRequestHandler)
	h := http.HandlerFunc(handler)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", pattern, nil)
	h.ServeHTTP(rr, req)

	app.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {

	// Do things

	txn, err := monitor.TxnFromRequest(w, r)
	if err == nil {
		txn.AddAttributes(log.Fields{
			"aString": "hello world",
			"aInt":    123,
		})
	}

	// Do more things

	if err == nil {
		txn.AddAttributes(log.Fields{
			"aString2": "goodbye",
			"aInt2":    456,
		})
	}
}
