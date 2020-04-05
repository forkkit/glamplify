package errors_test

import (
	"context"
	"errors"
	gerror "github.com/cultureamp/glamplify/errors"
	"github.com/cultureamp/glamplify/log"
	"testing"
	"time"
)

func TestHandler_Error_Success(t *testing.T) {

	gerror.HandleError(errors.New("NPE"), log.Fields{
		"user": "mike",
		"pwd": "abc",     // should be filtered out in bugsnag
		"age": 47,
	})

	time.Sleep(5 * time.Second)
}


func TestHandler_ErrorWithContext_Success(t *testing.T) {

	gerror.HandleErrorWithContext(context.TODO(), errors.New("NPE"), log.Fields{
		"user": "mike",
		"pwd": "abc",     // should be filtered out in bugsnag
		"age": 47,
	})

	time.Sleep(5 * time.Second)
}

