package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cultureamp/glamplify/config"
	"github.com/cultureamp/glamplify/event"
	"github.com/cultureamp/glamplify/log"
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
	// You can either get a new logger, or just use the public functions which internally use an internal logger
	// eg. log.Debug(), log.Print() and log.Error()

	// Example below shows usage with the package level logger (sensible default), but can
	// use an instance of a logger by calling log.New()

	// Emit debug trace
	// All messages must be static strings (as per Culture Amp Sensibile Default)
	log.Debug("Something happened")

	// Emit debug trace with fields
	// Fields can contain any type of variables
	log.Debug("Something happened", log.Fields{
		"aString": "hello",
		"aInt":    123,
		"aFloat":  42.48,
	})

	// Emit normal logging (can add optional fields if required)
	// Typically Print will be sent onto 3rd party aggregation tools (eg. Splunk)
	log.Print("Executing main")

	// Emit Error (can add optional fields if required)
	// Errors will always be sent onto 3rd party aggregation tools (eg. Splunk)
	err := errors.New("Main program stopped unexpectedly")
	log.Error(err)

	// If you want to set some fields for a particular scope (eg. for a Web Request
	// have a requestID for every log message within that scope) then you can use WithScope()
	scope := log.WithScope(log.Fields{"requestID": 123})

	// then just use the scope as you would a normal logger
	// Fields passed in the scope will be merged with any fields passed in subsequent calls
	// If duplicate keys, then fields in Debug, Print, Error will overwrite those of the scope
	scope.Print("Starting web request", log.Fields{"auth": "oauth"})

	// If you want to change the output or time format you can only do this for an
	// instance of the logger you create (not the internal one) by doing this:

	memBuffer := &bytes.Buffer{}
	logger := log.New(func(conf *log.Config) {
		conf.Output = memBuffer                 // can be set to anything that support io.Write
		conf.TimeFormat = "2006-01-02T15:04:05" // any valid time format
	})

	// The internall logger will always use these default values:
	// output = os.Stderr
	// time format = "2006-01-02T15:04:05.000Z07:00"
	logger.Print("Something useful just happened")

	/* EVENTS */

	app, err := event.NewApplication("GlamplifyUnitTests", func(conf *event.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
		conf.Labels = event.Labels{
			"asset":          "unknown",
			"classification": "restricted",
			"workload":       "development",
			"camp":           "amplify",
		}
	})

	_, handler := app.WrapHTTPHandler("/", rootRequestHandler)
	h := http.HandlerFunc(handler)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)

	app.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {

	// Do things

	txn, err := event.TxnFromRequest(w, r)
	if err == nil {
		txn.AddAttributes(event.Entries{
			"aString": "hello world",
			"aInt":    123,
		})
	}

	// Do more things

	if err == nil {
		txn.AddAttributes(event.Entries{
			"aString2": "goodbye",
			"aInt2":    456,
		})
	}
}
