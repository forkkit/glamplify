package log_test

import (
	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
	"net/http"
	"testing"
)

func Test_SeedRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "*", nil)

	req2 := log.SeedRequestWithRequestScopeFields(req)

	id, ok := log.GetTraceID(req2.Context())
	assert.Assert(t, ok && id != "", id)
}

func Test_RequestScope_AddGet(t *testing.T) {

	rsFields := log.RequestScopedFields{
		TraceID: "1-2-3",
		UserAggregateID: "a-b-c",
		CustomerAggregateID: "xyz",
	}

	req, _ := http.NewRequest("GET", "*", nil)

	req = log.AddRequestScopedFieldsRequest(req, rsFields)

	resultFields := log.GetRequestScopedFieldsRequest(req)
	assert.Assert(t, resultFields.TraceID == rsFields.TraceID, resultFields.TraceID)
	assert.Assert(t, resultFields.CustomerAggregateID == rsFields.CustomerAggregateID, resultFields.CustomerAggregateID)
	assert.Assert(t, resultFields.UserAggregateID == rsFields.UserAggregateID, resultFields.UserAggregateID)
}