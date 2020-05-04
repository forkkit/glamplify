package jwt

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
)

const (
	day  = 24 * time.Hour
	year = 365 * day // approx
)

func Test_JWT_Encode(t *testing.T) {

	jwt, err := NewEncoderFromPath("jwt.rs256.key.development.pem")
	assert.Assert(t, err == nil, err)

	expiry := 10 * year
	token, err := jwt.EncodeWithExpiry(Payload{
		Customer:      "abc123",
		RealUser:      "xyz234",
		EffectiveUser: "xyz345",
	}, expiry)

	assert.Assert(t, err == nil, err)
	splitToken := strings.Split(token, ".")
	assert.Assert(t, len(splitToken) == 3)

	header := splitToken[0]
	assert.Assert(t, header == "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9", header)

	fmt.Print(token)

	// Note: Hard to test payload/signature because there is a elapsed time field that always changes...
}

func Test_JWT_Decode(t *testing.T) {
	jwt, err := NewDecoderFromPath("jwt.rs256.key.development.pub")
	assert.Assert(t, err == nil, err)

	payload, err := jwt.Decode("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxNTg4NTY4OTY4LCJpYXQiOjE1ODg1NjgzNjh9.bQKvgB8ZfGSZDQnyfM1remW_cB_sF95iIS-QfJmdn3jOCK60xiwMB7cNkXBLeVsCyMScHiTyENePvCOudsMruNhWO8YvnBpv6712O4n2sRckTKaNdYXAwidILDFXzRZMzrFAOJu1zKjPSaXiXdEv6zgq3OMruXcBF9RzsOlKPoOlgBI1Q9ctGurgI-p4WovCA0YmV9w7I2c1t3WQMMlapoJPKW1-bM37sgEgJpJrmjCavYswQ_mWY0yk9h8ftXGvQRPvLXM_K-kkhmUJ1cLT-H4iXIZkCk-Y-ONAej9lPOgBGiCmOq5DHHcggOKzzqcT0YNKrZHfCrigd7ZbT-zRSw9ukzYafduabCSj9MAq_oKzYYbYpqu6yNtzHXFBZ7izWjGVMUxpQX5gaFh6W0aezWwmBL6drO1NzYDSMX2lJ-FwVCVfKbvqFPxS5mqYQCAQueGrTlrIndWqVdDbJFw2LHTFxVLAFQGgnM292WJYp6KYVKm07mRpzHdozb8ER7lfB_hlloudEBh14WxnV4iKZabjsGZmpzXldSdVKceXTBSY4jobE-vD_U2YfBcNU7y_A2qJtTnhdIWMq_UZrmi5ycV_Nq5MUSjLat-J8iFSkCeOEsyxQ3ybMiaxFEbpBZoZAIrRUJxx3KWtdzywyS4dN9frv36O0UuWomskaR1il6U")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, payload.Customer == "abc123", payload.Customer)
	assert.Assert(t, payload.RealUser == "xyz234", payload.RealUser)
	assert.Assert(t, payload.EffectiveUser == "xyz345", payload.EffectiveUser)
}
