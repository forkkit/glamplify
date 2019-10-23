package event

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/cultureamp/glamplify/log"

	newrelic "github.com/newrelic/go-agent"
)

// Entries contains key-value pairs to record along with the event
type Entries map[string]interface{}

// Labels are key value pairs used to roll up applications into specific categoriess
type Labels map[string]string

// Config contains Application and Transaction behavior settings.
// Use NewConfig to create a Config with proper defaults.
type Config struct {

	// Enabled controls whether the agent will communicate with the New Relic
	// servers and spawn goroutines.  Setting this to be false is useful in
	// testing and staging situations.
	Enabled bool

	// License is your New Relic license key.
	//
	// https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key
	License string

	// Logging controls whether Event logging is sent to StdOut or not
	Logging bool

	// Labels are key value pairs used to roll up applications into specific categories.
	//
	// https://docs.newrelic.com/docs/using-new-relic/user-interface-functions/organize-your-data/labels-categories-organize-apps-monitors
	Labels Labels

	// ServerlessMode contains fields which control behavior when running in
	// AWS Lambda.
	//
	// https://docs.newrelic.com/docs/serverless-function-monitoring/aws-lambda-monitoring/get-started/introduction-new-relic-monitoring-aws-lambda
	ServerlessMode bool

	// internal logger
	logger *eventLogger
}

// Application is a wrapper over the underlying implementation
type Application struct {
	impl newrelic.Application
	conf Config
}

// NewApplication creates a new Application - you should only create 1 Application per process
func NewApplication(name string, configure ...func(*Config)) (*Application, error) {

	conf := Config{
		Enabled:        false,
		Logging:        false,
		License:        os.Getenv("NEW_RELIC_LICENSE_KEY"),
		ServerlessMode: false,
		logger:         nil,
	}

	for _, config := range configure {
		config(&conf)
	}

	cfg := newrelic.NewConfig(name, conf.License)
	cfg.Enabled = conf.Enabled // useful to turn on/off in test/dev vs production accounts
	cfg.License = conf.License
	cfg.Labels = conf.Labels
	cfg.HighSecurity = false                // HighSecurity blocks sending custom events
	cfg.CustomInsightsEvents.Enabled = true // otherwise custom events won't fire
	cfg.Utilization.DetectAWS = true
	cfg.ServerlessMode.Enabled = conf.ServerlessMode
	cfg.ErrorCollector.Enabled = false
	cfg.ErrorCollector.CaptureEvents = false

	if conf.Logging {
		//cfg.Logger = newrelic.NewDebugLogger(os.Stdout) <- this writes JSON to Stdout :(
		// So we have our own implementation that wraps our standard logger

		conf.logger = newEventLogger()
		cfg.Logger = conf.logger
		cfg.ErrorCollector.Enabled = true
		cfg.ErrorCollector.CaptureEvents = true

		cfg.Logger.Debug("configuration", log.Fields{
			"enabled":        conf.Enabled,
			"logging":        conf.Logging,
			"labels":         conf.Labels,
			"ServerlessMode": conf.ServerlessMode,
		})
	}

	app := &Application{
		conf: conf,
	}

	impl, err := newrelic.NewApplication(cfg)
	if err != nil {
		app.logError("Failed to create Application", err)
		return nil, err
	}

	if !conf.ServerlessMode {
		// if conf.ServerlessMode = false (server mode) then newrelic.NewApplication spins up
		// some go routines that make a network call back to NR. Until this happens any "RecordCustomEvents"
		// seem to get dropped!
		// Waiting here so that everything is set up and ready
		time.Sleep(5 * time.Second)
	}

	app.impl = impl
	return app, err
}

// RecordEvent sends a custom event with the associated data to the underlying implementation
func (app Application) RecordEvent(eventType string, entries Entries) error {
	app.log("Begin RecordEvent",
		log.Fields{
			"eventType": eventType,
		},
		entries,
	)
	err := app.impl.RecordCustomEvent(eventType, entries)
	app.logError("RecordEvent", err)

	return err
}

// WrapHTTPHandler adds a Transaction within the current request
func (app *Application) WrapHTTPHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	p, h := app.wrapHTTPHandler(pattern, http.HandlerFunc(handler))
	return p, func(w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) }
}

// Shutdown flushes any remaining data to the SAAS endpoint
func (app Application) Shutdown() {

	if !app.conf.ServerlessMode {
		// if conf.ServerlessMode = false (server mode) then newrelic.Shutdown can exit its internal go routines
		// before it has sent all pending data!
		// Waiting here so that everything is sent before we start closing down...
		time.Sleep(5 * time.Second)
	}

	// The time duration passed here is how long to wait before the shutdown channel processes the request
	// It is NOT how long to wait to send data before shutting down.
	app.impl.Shutdown(30 * time.Second)
}

func (app *Application) wrapHTTPHandler(pattern string, handler http.Handler) (string, http.Handler) {
	return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := app.startTransaction(pattern, w, r)
		defer txn.End()

		handler.ServeHTTP(txn, r)
	})
}

func (app *Application) startTransaction(name string, w http.ResponseWriter, r *http.Request) Transaction {
	app.log("Starting Transaction", log.Fields{
		"txnName": name,
	})

	// Create our wrapper txn and add it to the HTTP ctx
	txn := Transaction{
		app:     app,
		name:    name,
		logging: app.conf.Logging,
		logger:  app.conf.logger,
	}
	r = txn.addToHTTPContext(r)

	// call the NR implementation
	impl := app.impl.StartTransaction(name, w, r)
	txn.impl = impl

	return txn
}

func (app *Application) addToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, appContextKey, app)
}

func (app Application) log(msg string, fields log.Fields, entries ...Entries) {
	if app.conf.Logging {
		merged := app.conf.logger.merge(fields, entries...)
		app.conf.logger.Debug(msg, merged)
	}
}

func (app Application) logError(msg string, err error) {
	if err != nil && app.conf.Logging {
		app.conf.logger.Error(msg, map[string]interface{}{
			"error": err,
		})
	}
}
