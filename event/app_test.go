package event_test

import (
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

func TestApplication_RecordEvent_Server_Fail(t *testing.T) {
	app, err := event.NewApplication("Glamplify-Unit-Tests", func(conf *event.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	assert.Assert(t, err == nil, err)
	assert.Assert(t, app != nil, "application was nil")

	err = app.RecordEvent("glamplify_unittest_customevent", event.Entries{
		"test":  "big_long_string_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890",
		"anull": nil,
	})
	assert.Assert(t, err != nil, err)

	app.Shutdown()
}

