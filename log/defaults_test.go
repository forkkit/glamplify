package log

import (
	"fmt"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_TraceID(t *testing.T) {

	id := traceID()
	fmt.Println(id)

	assert.Assert(t, strings.Contains(id, "-"))
	assert.Assert(t, strings.Count(id, "-") == 2)
}
