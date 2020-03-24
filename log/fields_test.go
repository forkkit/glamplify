package log_test

import (
	"github.com/cultureamp/glamplify/log"
	"testing"

	"gotest.tools/assert"
)

func TestEntries_Success(t *testing.T) {
	entries := log.Fields{
		"aString": "hello world",
		"aInt":    123,
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Assert(t, ok, ok)
	assert.Assert(t, err == nil, err)
}

func TestEntries_InvalidType_Failed(t *testing.T) {
	dict := map[string]int{
		"key1": 1,
	}
	entries := log.Fields{
		"aMap": dict,
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}

func TestEntries_NilValue_Failed(t *testing.T) {
	dict := map[string]interface{}{
		"key1": nil,
	}
	entries := log.Fields{
		"aMap": dict,
		"akey": nil,
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}

func TestEntries_StringToLong_Failed(t *testing.T) {
	entries := log.Fields{
		"aString": "big_long_string_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890",
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}
