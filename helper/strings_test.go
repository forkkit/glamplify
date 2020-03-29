package helper

import (
	"gotest.tools/assert"
	"testing"
)

func Test_ToSnakeCase(t *testing.T) {

	sc := ToSnakeCase("hello")
	assert.Assert(t,  sc == "hello", "was: '%s'", sc)

	sc = ToSnakeCase("requestID")
	assert.Assert(t,  sc == "request_id", "was: '%s'", sc)

	sc = ToSnakeCase("LastGC")
	assert.Assert(t,  sc == "last_gc", "was: '%s'", sc)

	sc = ToSnakeCase("request_id")
	assert.Assert(t,  sc == "request_id", "was: '%s'", sc)

	sc = ToSnakeCase("something happened")
	assert.Assert(t,  sc == "something_happened", "was: '%s'", sc)

	sc = ToSnakeCase(" with  added  spaces")
	assert.Assert(t,  sc == "with_added_spaces", "was: '%s'", sc)

	sc = ToSnakeCase(" And  WITH  Capitals  ")
	assert.Assert(t,  sc == "and_with_capitals","was: '%s'", sc)
}

func BenchmarkLogging(b *testing.B) {

	sa := []string {"hello", "requestID", "request_id", "something happened"}

	for n := 0; n < b.N; n++ {
		for _, s := range sa {
			ToSnakeCase(s)
		}
	}
}
