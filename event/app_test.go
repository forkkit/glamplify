package event_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cultureamp/glamplify/event"
	"gotest.tools/assert"
)

func TestApplication_RecordEvent_Server_Success(t *testing.T) {
	app, err := event.NewApplication("Glamplify-Unit-Tests", func(conf *event.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	assert.Assert(t, err == nil, err)
	assert.Assert(t, app != nil, "application was nil")

	err = app.RecordEvent("glamplify_unittest_customevent", event.Entries{
		"aString": "hello world",
		"aInt":    123,
	})
	assert.Assert(t, err == nil, err)

	app.Shutdown()
}

func TestApplication_AddAttribute_Server_Success(t *testing.T) {
	app, err := event.NewApplication("Glamplify-Unit-Tests", func(conf *event.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	assert.Assert(t, err == nil, err)
	assert.Assert(t, app != nil, "application was nil")

	_, handler := app.WrapTxnHandler("/", addAttribute)
	h := http.HandlerFunc(handler)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)

	app.Shutdown()
}

func addAttribute(w http.ResponseWriter, r *http.Request) {
	txn, ok := event.TxnFromRequest(w, r)
	if ok {
		txn.AddAttribute("txnString", "hello txn world")
		txn.AddAttribute("txnInt", "456")
	}
}
