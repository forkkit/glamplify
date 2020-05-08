package aws

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"gotest.tools/assert"
	"testing"
)

func Test_GetSecretParam_MissingKey(t *testing.T) {

	sm := NewSecretsManager("default")
	assert.Assert(t, sm != nil, sm)

	// Missing Key
	val, err := sm.Get("/this/should/not/exist/secret_key")
	assert.Assert(t, val == "", val)

	aerr, ok := err.(awserr.Error)
	assert.Assert(t, ok, ok)
	assert.Assert(t, aerr.Message() != "", aerr.Message())
}

// TODO - what is a good key to use for unit tests?

