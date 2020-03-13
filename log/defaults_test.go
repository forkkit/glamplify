package log

import (
	"fmt"
	"github.com/cultureamp/glamplify/constants"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_HostName(t *testing.T) {

	df := NewDefaultValues(constants.RFC3339Milli)
	host := df.hostName()

	assert.Assert(t, host != "", host)
	assert.Assert(t, host != "<unknown>", host)
}

func Test_TraceID(t *testing.T) {

	df := NewDefaultValues(constants.RFC3339Milli)
	id := df.NewTraceID()
	fmt.Println(id)

	assert.Assert(t, strings.Contains(id, "-"))
	assert.Assert(t, strings.Count(id, "-") == 2)
}
