package event

import (
	"os"

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

	// Labels are key value pairs used to roll up applications into specific categories.
	//
	// https://docs.newrelic.com/docs/using-new-relic/user-interface-functions/organize-your-data/labels-categories-organize-apps-monitors
	Labels Labels

	// ServerlessMode contains fields which control behavior when running in
	// AWS Lambda.
	//
	// https://docs.newrelic.com/docs/serverless-function-monitoring/aws-lambda-monitoring/get-started/introduction-new-relic-monitoring-aws-lambda
	ServerlessMode bool
}

// Application is a wrapper over the underlying implementation
type Application struct {
	nrapp newrelic.Application
}

// NewApplication creates a new Application - you should only create 1 Application per process
func NewApplication(name string, configure ...func(*Config)) (*Application, error) {
	app := &Application{}

	lic := os.Getenv("NEW_RELIC_LICENSE_KEY")
	conf := Config{
		Enabled:        true,
		ServerlessMode: false,
	}

	for _, config := range configure {
		config(&conf)
	}

	cfg := newrelic.NewConfig(name, lic)
	//cfg.Logger = newrelic.NewDebugLogger(os.Stdout) <- this writes JSON to Stdout :( Need our own logger impl?
	cfg.Enabled = conf.Enabled
	cfg.Labels = conf.Labels
	cfg.CustomInsightsEvents.Enabled = true
	cfg.Utilization.DetectAWS = true
	cfg.ServerlessMode.Enabled = conf.ServerlessMode

	nrapp, err := newrelic.NewApplication(cfg)
	app.nrapp = nrapp
	return app, err
}

// RecordEvent sends a custom event with the associated data to the underlying implementation
func (app Application) RecordEvent(name string, entries Entries) {
	app.nrapp.RecordCustomEvent(name, entries)
}
