package monitor_test

import (
	"github.com/cultureamp/glamplify/monitor"
	"testing"

	"gotest.tools/assert"
)

func TestEntries_Success(t *testing.T) {
	entries := monitor.Fields{
		"aString": "hello world",
		"aInt":    123,
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.Validate()
	assert.Assert(t, ok, ok)
	assert.Assert(t, err == nil, err)
}

func TestEntries_InvalidType_Failed(t *testing.T) {
	dict := map[string]int{
		"key1": 1,
	}
	entries := monitor.Fields{
		"aMap": dict,
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.Validate()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}

func TestEntries_StringToLong_Failed(t *testing.T) {
	entries := monitor.Fields{
		"aString": "big_long_string_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890",
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.Validate()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}
