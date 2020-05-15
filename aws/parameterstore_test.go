package aws

import (
	"gotest.tools/assert"
	"testing"
)

func Test_GetParam_MissingKey(t *testing.T) {

	ps := NewParameterStore("default")
	assert.Assert(t, ps != nil, ps)

	// Missing Key
	val, err := ps.Get("/this/should/not/exist/secret_key")
	assert.Assert(t, val == "", val)
	assert.Assert(t, err != nil, val)

	// aerr, ok := err.(awserr.Error)
}

// TODO - what is a good key to use for unit tests?
