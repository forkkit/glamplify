package monitor_test

import (
	"github.com/cultureamp/glamplify/monitor"
	"testing"

	"gotest.tools/assert"
)

func TestApplication_RecordEvent_Server_Success(t *testing.T) {
	app, err := monitor.NewApplication("Glamplify-Unit-Tests", func(conf *monitor.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	assert.Assert(t, err == nil, err)
	assert.Assert(t, app != nil, "application was nil")

	err = app.RecordEvent("glamplify_unittest_customevent", monitor.Entries{
		"aString": "hello world",
		"aInt":    123,
	})
	assert.Assert(t, err == nil, err)

	app.Shutdown()
}

func TestApplication_RecordEvent_Server_Fail(t *testing.T) {
	app, err := monitor.NewApplication("Glamplify-Unit-Tests", func(conf *monitor.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	assert.Assert(t, err == nil, err)
	assert.Assert(t, app != nil, "application was nil")

	err = app.RecordEvent("glamplify_unittest_customevent", monitor.Entries{
		"test":  "big_long_string_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890",
		"anull": nil,
	})
	assert.Assert(t, err != nil, err)

	app.Shutdown()
}

