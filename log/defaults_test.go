package log

import (
	"gotest.tools/assert"
	"testing"
)

func Test_HostName(t *testing.T) {

	host := hostName()

	assert.Assert(t, host != "", host)
	assert.Assert(t, host != "<unknown>", host)
}
