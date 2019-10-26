package event_test

import (
	"testing"

	"github.com/cultureamp/glamplify/event"
	"gotest.tools/assert"
)

func TestEntries_Success(t *testing.T) {
	entries := event.Entries{
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
	entries := event.Entries{
		"aMap": dict,
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.Validate()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}

func TestEntries_StringToLong_Failed(t *testing.T) {
	entries := event.Entries{
		"aString": "big_long_string_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890",
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.Validate()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}
