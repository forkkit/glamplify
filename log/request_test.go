package log_test

import (
	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
	"net/http"
	"testing"
)

func Test_SeedRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "*", nil)

	req2 := log.SeedRequestCtxWithRequestScopeFields(req)

	id, ok := log.GetTraceID(req2.Context())
	assert.Assert(t, ok && id != "", id)
}
