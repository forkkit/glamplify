package event_test

import (
	"github.com/cultureamp/glamplify/event"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestApplication_RecordEvent_Server_Success(t *testing.T) {
	app, err := event.NewApplication("Glamplify-Unit-Tests", func(conf *event.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	assert.Assert(t, err == nil, err)
	assert.Assert(t, app != nil, "application was nil")

	time.Sleep(3 * time.Second)

	err = app.RecordEvent("glamplify_unittest_customevent", event.Entries{
		"aString": "hello world",
		"aInt": 123,
	})
	assert.Assert(t, err == nil, err)

	time.Sleep(3 * time.Second)
	app.Shutdown()
}

