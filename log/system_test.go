package log

import (
	"errors"
	"gotest.tools/assert"
	"testing"
)

func Test_HostName(t *testing.T) {

	df := newSystemValues()
	host := df.hostName()

	assert.Assert(t, host != "", host)
	assert.Assert(t, host != "<unknown>", host)
}

func Test_Default(t *testing.T) {
	df := newSystemValues()

	fields := df.getSystemValues(rsFields, "event_name", DebugSev)

	_, ok := fields[Time]
	assert.Assert(t, ok, "missing 'time' in default fields")
	_, ok = fields[Event]
	assert.Assert(t, ok, "missing 'event' in default fields")
	_, ok = fields[Resource]
	assert.Assert(t, ok, "missing 'resource' in default fields")
	_, ok = fields[Os]
	assert.Assert(t, ok, "missing 'os' in default fields")
	_, ok = fields[Severity]
	assert.Assert(t, ok, "missing 'severity' in default fields")

	_, ok = fields[TraceID]
	assert.Assert(t, ok, "missing 'trace_id' in default fields")
	_, ok = fields[Customer]
	assert.Assert(t, ok, "missing 'customer' in default fields")
	_, ok = fields[User]
	assert.Assert(t, ok, "missing 'user' in default fields")

	_, ok = fields[Product]
	assert.Assert(t, ok, "missing 'product' in default fields")
	_, ok = fields[App]
	assert.Assert(t, ok, "missing 'app' in default fields")
	_, ok = fields[AppVer]
	assert.Assert(t, ok, "missing 'app_ver' in default fields")
	_, ok = fields[AwsRegion]
	assert.Assert(t, ok, "missing 'region' in default fields")
}

func Test_ErrorDefault(t *testing.T) {
	df := newSystemValues()

	fields := df.getSystemValues(rsFields, "event_name", DebugSev)
	fields = df.getErrorValues(errors.New("test err"), fields)

	_, ok := fields[Exception]
	assert.Assert(t, ok, "missing 'exception' in default fields")
}