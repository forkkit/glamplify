package helper

import (
	"fmt"
	"gotest.tools/assert"
	"strings"
	"testing"
	"time"
)

func Test_TraceID(t *testing.T) {

	id := NewTraceID()
	fmt.Println(id)

	assert.Assert(t, strings.Contains(id, "-"))
	assert.Assert(t, strings.Count(id, "-") == 2)
}


func Test_DurationAsIso8601(t *testing.T) {
	d := time.Millisecond * 456
	s := DurationAsISO8601(d)
	assert.Assert(t, s == "P0.456S", "was: %s", s)

	d = time.Millisecond * 1456
	s = DurationAsISO8601(d)
	assert.Assert(t, s == "P1.456S", "was: %s", s)

}
