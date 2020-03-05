package http

import (
	"github.com/cultureamp/glamplify/monitor"
	"github.com/cultureamp/glamplify/notify"
	"net/http"
)

func WrapHTTPHandler(
	app *monitor.Application,
	notify *notify.Notifier,
	pattern string,
	handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {

	// 1. Wrap with bugsnag
	pattern, handler = notify.WrapHTTPHandler(pattern, handler)

	// 2. Then wrap with new relic
	return app.WrapHTTPHandler(pattern, handler)
}
