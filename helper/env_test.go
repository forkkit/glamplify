package helper

import (
	"gotest.tools/assert"
	"testing"
)

func Test_EnvDefault(t *testing.T) {

	val := GetEnvOrDefault("should_not_exist_env_var", "fallback")
	assert.Assert(t, val == "fallback", val)
}

func Test_EnvGet(t *testing.T) {

	val := GetEnvOrDefault("PATH", "fallback")
	assert.Assert(t, val != "fallback", val)
}
