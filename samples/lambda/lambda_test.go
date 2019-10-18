package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/djhworld/go-lambda-invoke/golambdainvoke"
	"gotest.tools/assert"
)

func TestApplication_Lambda_Success(t *testing.T) {

	response, err := golambdainvoke.Run(golambdainvoke.Input{
		Port:    8001,
		Payload: "test",
	})

	assert.Assert(t, err == nil, err)

	var result string
	if err = json.Unmarshal(response, &result); err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	assert.Assert(t, result == "Ok", result)
}
