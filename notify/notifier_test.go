package notify_test


import (
	"errors"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/notify"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestNotify_Error_Success(t *testing.T) {

	err := notify.Error(errors.New("NPE"), log.Fields{
		"user": "mike",
		"pwd": "abc",     // should be filtered out in bugsnag
		"age": 47,
	})

	assert.Assert(t, err == nil, err)
	time.Sleep(5 * time.Second)
}

