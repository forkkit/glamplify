package jwt

import (
	"fmt"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_JWT_Encode(t *testing.T) {

	jwt, err := NewJWTEncoderFromPath("jwt.rs256.key.development.pem")
	assert.Assert(t, err == nil, err)

	token, err := jwt.Encode(Payload{
		Customer:      "abc123",
		RealUser:      "xyz234",
		EffectiveUser: "xyz345",
	})
	assert.Assert(t, err == nil, err)
	splitToken := strings.Split(token, ".")
	assert.Assert(t, len(splitToken) == 3)

	header := splitToken[0]
	assert.Assert(t, header == "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9", header)

	fmt.Print(token)

	// Note: Hard to test payload/signature because there is a elapsed time field that always changes...
}

func Test_JWT_Decode(t *testing.T) {

	jwt, err := NewJWTDecoderFromPath("jwt.rs256.key.development.pub")
	assert.Assert(t, err == nil, err)

	payload, err := jwt.Decode("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxNTg1NzkzMDk1LCJpYXQiOjE1ODU3OTI0OTV9.AURdMTjJQwx4cUeuH_MiiBvIxQjUCAogWFELbTcOyi2oY5mGY3TY9kGjyxYHy2GS-2AQD7Kh7FHEaoWAsqn0c7YLsZWu2bYBbvlCs6Rxi_RX9SJ-RM8WbsW87dJYtQxrbZHDtZXAJXjjuUbBkVEaKH3hYtePSBDINpFPBrpFIBc8q0VBfaGf01r0A32bvd9ztOAQKA5nE6ZsgehKL8oPVc76z0-CykrRjjI4YFVuMiooXqMkCYjKe_5IN8fVVSS5GuOrVeyfqDZlDAVs9oOYKx9lXzqI0Vgvnu5Yd1loctoLtfZAcXHJA1wit033UCrlGjGqAW6E9lUY87Y65Ip4lLatEV2OW8E0_sGciovyFsd3bD06XbltQo4AxaKNcKfv4pE5-wffx7D5Pg4hPEptyicIr8Q37x4yORV7IyGGuuo7TSbxLjZx8SYAeCKsBArw0DTocB86NiuIJ3bKaJK3lFqu1OpQW_UXmyIRwvf0PGkN6PFtHmu2--dedtdaW-0vSKt69ie93Vv5JpFCy26WEvy4P2LaxMS-m8WBRa5L9RIpDY4UTRpFwoV0q5z7ae3Cx4zdkFkQDD4EIzMKdLM6m-6WjI5N8_7hxgvLc4u1rU0uWWxMN6bhpsv2As1I2S94IKjqL70q2SnjJnQhB5x65N7VJWKU0gFOR847yjYO4k0")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, payload.Customer == "abc123", payload.Customer)
	assert.Assert(t, payload.RealUser == "xyz234", payload.RealUser)
	assert.Assert(t, payload.EffectiveUser == "xyz345", payload.EffectiveUser)
}