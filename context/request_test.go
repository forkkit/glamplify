package context_test

import (
	"github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/jwt"
	"gotest.tools/assert"
	"net/http"
	"testing"
)


func Test_RequestScope_AddGet(t *testing.T) {

	rsFields := context.RequestScopedFields{
		TraceID: "1-2-3",
		UserAggregateID: "a-b-c",
		CustomerAggregateID: "xyz",
	}

	req, _ := http.NewRequest("GET", "*", nil)
	req = context.AddRequestScopedFieldsRequest(req, rsFields)

	resultFields, ok := context.GetRequestScopedFieldsFromRequest(req)
	assert.Assert(t, ok, ok)
	assert.Assert(t, resultFields.TraceID == rsFields.TraceID, resultFields.TraceID)
	assert.Assert(t, resultFields.CustomerAggregateID == rsFields.CustomerAggregateID, resultFields.CustomerAggregateID)
	assert.Assert(t, resultFields.UserAggregateID == rsFields.UserAggregateID, resultFields.UserAggregateID)
}

func Test_Request_Wrap(t *testing.T) {
	rsFields := context.RequestScopedFields{
		TraceID: "1-2-3",
		RequestID: "7-8-9",
		UserAggregateID: "a-b-c",
		CustomerAggregateID: "xyz",
	}

	req, _ := http.NewRequest("GET", "*", nil)
	req = context.AddRequestScopedFieldsRequest(req, rsFields)

	req = context.WrapRequest(req)
	id, ok := context.GetTraceID(req.Context())
	assert.Assert(t, ok && id == "1-2-3", id)
	id, ok = context.GetRequestID(req.Context())
	assert.Assert(t, ok && id == "7-8-9", id)
	id, ok = context.GetUser(req.Context())
	assert.Assert(t, ok && id == "a-b-c", id)
	id, ok = context.GetCustomer(req.Context())
	assert.Assert(t, ok && id == "xyz", id)
}

func Test_Request_WrapWithDecoder(t *testing.T) {
	jwt, err := jwt.NewDecoderFromPath("../jwt/jwt.rs256.key.development.pub")
	assert.Assert(t, err == nil, err)

	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxOTAzOTMwNzA0LCJpYXQiOjE1ODg1NzA3MDR9.XGm34FDIgtBFvx5yC2HTUu-cf3DaQI4TmIBVLx0H7y89oNVNWJaKA3dLvWS0oOZoYIuGhj6GzPREBEmou2f9JsUerqnc-_Tf8oekFZWU7kEfzu9ECBiSWPk7ljPJeZLbau62sSqD7rYb-m3v1mohqz4tKJ_7leWu9L1uHHliC7YGlSRl1ptVDllJjKXKjOg9ifeGSXDEMeU35KgCFwIwKdu8WmCTd8ztLSKEnLT1OSaRZ7MSpmHQ4wUZtS6qvhLBiquvHub9KdQmc4mYWLmfKdDiR5DH-aswJFGLVu3yisFRY8uSfeTPQRhQXd_UfdgifCTXdWTnCvNZT-BxULYG-5mlvAFu-JInTga_9-r-wHRzFD1SrcKjuECF7vUG8czxGNE4sPjFrGVyBxE6fzzcFsdrhdqS-LB_shVoG940fD-ecAhXQZ9VKgr-rmCvmxuv5vYI2HoMfg9j_-zeXkucKxvPYvDQZYMdeW4wFsUORliGplThoHEeRQxTX8d_gvZFCy_gGg0H57FmJwCRymWk9v29s6uyHUMor_r-e7e6ZlShFBrCPAghXL04S9IFJUxUv30wNie8aaSyvPuiTqCgGiEwF_20ZaHCgYX0zupdGm4pHTyJrx2wv31yZ4VZYt8tKjEW6-BlB0nxzLGk5OUN83vq-RzH-92WmY5kMndF6Jo"

	req, _ := http.NewRequest("GET", "*", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	req2 := context.WrapRequestWithDecoder(req, jwt)
	id, ok := context.GetTraceID(req2.Context())
	assert.Assert(t, ok && id != "", id)
	id, ok = context.GetCustomer(req2.Context())
	assert.Assert(t, ok && id != "", id)
	id, ok = context.GetUser(req2.Context())
	assert.Assert(t, ok && id != "", id)
}
