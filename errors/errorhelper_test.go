package errors_test

import (
	"context"
	"errors"
	"github.com/cultureamp/glamplify/types"
	"testing"
	"time"
)

func TestHandler_Error_Success(t *testing.T) {

	HandleError(errors.New("NPE"), types.Fields{
		"user": "mike",
		"pwd": "abc",     // should be filtered out in bugsnag
		"age": 47,
	})

	time.Sleep(5 * time.Second)
}


func TestHandler_ErrorWithContext_Success(t *testing.T) {

	HandleErrorWithContext(errors.New("NPE"), context.TODO(), types.Fields{
		"user": "mike",
		"pwd": "abc",     // should be filtered out in bugsnag
		"age": 47,
	})

	time.Sleep(5 * time.Second)
}

