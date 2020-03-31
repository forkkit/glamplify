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
	assert.Assert(t, token == "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0In0.oDXzd1tq5XRpENEcw7GxAglOpRWmL5ld7XYPeNlrF-IfWYYRy86rta9yG9ug5wS1GV7Lvv8EbufXk0DKTnd23oObWoJtXLUaHh2TG9sw9bsxNwLKu1eWw7MQtUYByN2QFpRGeMQo_yw5Y6bT76janQ1NZknopHHvttcLBFuSMdThMX-4gOlaCuVsr8MQ218WUC-rVrSAol57at_2gf8PkEik3bcOd4bvUpf-ThumkljyzSrxVBY57H1kYbYAST4CwcCrf2F3oTLa_xNbFycngVCvJLZtSQR5GxwpO_ERqFziEaQ07bW6Svcs0EvARjCB-4vYdKTFaw3J5qu2aWVHf9m3a4QPA5O91ODYFYq_7k6upmxQl074_MQ-ZsnDRt0cUyPJjObMjU99MuMLQNnAMU67iNYkOxocR1OCNzLL1ObpeoYVq8sZWQPVhrPFDnC-V5uIsoSl9NofwcApLfUV2WjcMHxPfJYqPo-BNq3P_p1G1WSJ7iLP1BMXAU_ZaK49YaWb3fwu4NzRSCjsulWjMiE1yQL_bQrj4crygAyCgG7hpgq9OdiVl7YElrOL-oY1_3XCvnVcZkCd5dQjSbTXx-cW8Xc_zeY1QGxiKaeI3Yg24XLSVSFMNX4XNXwtNlK-LSrWQU8S0bVZBRDNo0jM9hx7INjYc4tamu2sGcH-71Q", token)
}

func Test_JWT_Decode(t *testing.T) {
	ctx := context.Background()
	jwt := NewJWTDecoderFromPath(ctx, "jwt.rs256.key.development.pub")

	payload, err := jwt.Decode("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0In0.oDXzd1tq5XRpENEcw7GxAglOpRWmL5ld7XYPeNlrF-IfWYYRy86rta9yG9ug5wS1GV7Lvv8EbufXk0DKTnd23oObWoJtXLUaHh2TG9sw9bsxNwLKu1eWw7MQtUYByN2QFpRGeMQo_yw5Y6bT76janQ1NZknopHHvttcLBFuSMdThMX-4gOlaCuVsr8MQ218WUC-rVrSAol57at_2gf8PkEik3bcOd4bvUpf-ThumkljyzSrxVBY57H1kYbYAST4CwcCrf2F3oTLa_xNbFycngVCvJLZtSQR5GxwpO_ERqFziEaQ07bW6Svcs0EvARjCB-4vYdKTFaw3J5qu2aWVHf9m3a4QPA5O91ODYFYq_7k6upmxQl074_MQ-ZsnDRt0cUyPJjObMjU99MuMLQNnAMU67iNYkOxocR1OCNzLL1ObpeoYVq8sZWQPVhrPFDnC-V5uIsoSl9NofwcApLfUV2WjcMHxPfJYqPo-BNq3P_p1G1WSJ7iLP1BMXAU_ZaK49YaWb3fwu4NzRSCjsulWjMiE1yQL_bQrj4crygAyCgG7hpgq9OdiVl7YElrOL-oY1_3XCvnVcZkCd5dQjSbTXx-cW8Xc_zeY1QGxiKaeI3Yg24XLSVSFMNX4XNXwtNlK-LSrWQU8S0bVZBRDNo0jM9hx7INjYc4tamu2sGcH-71Q")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, payload.Customer == "abc123", payload.Customer)
	assert.Assert(t, payload.RealUser == "xyz234", payload.RealUser)
	assert.Assert(t, payload.EffectiveUser == "xyz345", payload.EffectiveUser)
}