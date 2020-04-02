package aws

import (
	"gotest.tools/assert"
	"testing"
)

func Test_GetParam(t *testing.T) {

	ps := NewParameterStore("default")
	assert.Assert(t, ps != nil, ps)

	// TODO - what is a good key to use for unit tests?
	// Should I mock this?
	// _, err := ps.Get("common/AUTH_PUBLIC_KEY")
	// assert.Assert(t, err == nil, err)
}
