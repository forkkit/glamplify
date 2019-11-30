package notify_test

import (
	"context"
	"errors"
	"github.com/cultureamp/glamplify/field"
	"github.com/cultureamp/glamplify/notify"
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotify_Error_Success(t *testing.T) {

	notifier, err := notify.NewNotifier("GlamplifyUnitTests", func (conf *notify.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	assert.Assert(t, err == nil, err)

	err = notifier.Error(errors.New("NPE"), field.Fields{
		"user": "mike",
		"pwd": "abc",     // should be filtered out in bugsnag
		"age": 47,
	})
	assert.Assert(t, err == nil, err)

	notifier.Shutdown()
}


func TestNotify_Context_Success(t *testing.T) {

	notifier, err := notify.NewNotifier("GlamplifyUnitTests", func (conf *notify.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	assert.Assert(t, err == nil, err)

	_, handler := notifier.WrapHTTPHandler("/", rootRequest)
	h := http.HandlerFunc(handler)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)

	// Add *testing.T to request context
	ctx := req.Context()
	ctx = context.WithValue(ctx, "t", t)
	req = req.WithContext(ctx)

	h.ServeHTTP(rr, req)

	notifier.Shutdown()
}


func rootRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	t, _ := ctx.Value("t").(*testing.T)

	notifier, err := notify.NotifyFromContext(ctx)
	assert.Assert(t, err == nil, err)

	err = notifier.ErrorWithContext(errors.New("NPE"), ctx, field.Fields{
		"user": "mike",
		"pwd": "abc",     // should be filtered out in bugsnag
		"age": 47,
	})
	assert.Assert(t, err == nil, err)
}


