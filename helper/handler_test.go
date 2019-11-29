package helper_test

import (
	"context"
	"errors"
	"github.com/cultureamp/glamplify/field"
	"github.com/cultureamp/glamplify/helper"
	"testing"
	"time"
)

func TestHandler_Error_Success(t *testing.T) {

	helper.HandleError(errors.New("NPE"), field.Fields{
		"user": "mike",
		"pwd": "abc",     // should be filtered out in bugsnag
		"age": 47,
	})

	time.Sleep(5 * time.Second)
}


func TestHandler_ErrorWithContext_Success(t *testing.T) {

	helper.HandleErrorWithContext(errors.New("NPE"), context.TODO(), field.Fields{
		"user": "mike",
		"pwd": "abc",     // should be filtered out in bugsnag
		"age": 47,
	})

	time.Sleep(5 * time.Second)
}

