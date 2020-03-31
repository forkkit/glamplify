package jwt


import (
	"context"
	"gotest.tools/assert"
	"testing"
)

func Test_JWT_Encode(t *testing.T) {
	ctx := context.Background()
	jwt := NewJWTEncoderFromPath(ctx, "jwt.rs256.key.development.pem")


	token, err := jwt.Encode(Payload{
		Customer:      "abc123",
		RealUser:      "xyz234",
		EffectiveUser: "xyz345",
	})
	assert.Assert(t, err == nil, err)
	assert.Assert(t, token == "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0In0.a4MOmuPZBfg8rs0Mz25HA6yHkz5wW7V9whI6wX-0KCspSBNbPK2ocdjqtLZ0slcQcppeCh32W0l-AOoyL3_IGX5r8U3ldCw_hdgOAGfb0h8Kk4LhUOfOLOOG_NwwnpyJfZ0sqO5UHyuFwW3vnJVEJBLsJIMPEAFZyqxccuZzZhtSJqjGIIVB9nSzt9SE2Tpeu6lejvYyNrAj-aPYv_lnt0MAXV_HUUmaaSYnhIUYhc0IQPEvA37Wk9cY6heM9BpnVbcqtw8F1WyXYZlOBRBVguKtbATVNKlnSrqVifsPixTr0-bebFWta4lptSPjO8DogN5YjPgyyI6tnJlOfhsWHQ", token)
}

func Test_JWT_Decode(t *testing.T) {
	ctx := context.Background()
	jwt := NewJWTDecoderFromPath(ctx, "jwt.rs256.key.development.pub")

	payload, err := jwt.Decode("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0In0.a4MOmuPZBfg8rs0Mz25HA6yHkz5wW7V9whI6wX-0KCspSBNbPK2ocdjqtLZ0slcQcppeCh32W0l-AOoyL3_IGX5r8U3ldCw_hdgOAGfb0h8Kk4LhUOfOLOOG_NwwnpyJfZ0sqO5UHyuFwW3vnJVEJBLsJIMPEAFZyqxccuZzZhtSJqjGIIVB9nSzt9SE2Tpeu6lejvYyNrAj-aPYv_lnt0MAXV_HUUmaaSYnhIUYhc0IQPEvA37Wk9cY6heM9BpnVbcqtw8F1WyXYZlOBRBVguKtbATVNKlnSrqVifsPixTr0-bebFWta4lptSPjO8DogN5YjPgyyI6tnJlOfhsWHQ"	)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, payload.Customer == "abc123", payload.Customer)
	assert.Assert(t, payload.RealUser == "xyz234", payload.RealUser)
	assert.Assert(t, payload.EffectiveUser == "xyz345", payload.EffectiveUser)
}