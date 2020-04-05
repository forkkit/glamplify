package log

import (
	"fmt"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_HostName(t *testing.T) {

	df := newDefaultValues()
	host := df.hostName()

	assert.Assert(t, host != "", host)
	assert.Assert(t, host != "<unknown>", host)
}

func Test_TraceID(t *testing.T) {

	df := newDefaultValues()
	id := df.newTraceID()
	fmt.Println(id)

	assert.Assert(t, strings.Contains(id, "-"))
	assert.Assert(t, strings.Count(id, "-") == 2)
}
