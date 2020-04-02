package http

import (
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/monitor"
	"github.com/cultureamp/glamplify/notify"
	"gotest.tools/assert"
	"net/http"
	"testing"
)

func Test_Wrap(t *testing.T) {

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
	assert.Assert(t, appErr == nil, appErr)

	notifier, notifyErr := notify.NewNotifier("GlamplifyUnitTests", func (conf *notify.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	assert.Assert(t, notifyErr == nil, notifyErr)

	pattern, handler := WrapHTTPHandler(app, notifier, "/", rootRequestHandler)
	assert.Assert(t, handler != nil, handler)
	assert.Assert(t, pattern == "/", pattern)

}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {}