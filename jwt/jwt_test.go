package jwt

import (
	"fmt"
	"gotest.tools/assert"
	"strings"
	"testing"
	"time"
)

const (
	day  = 24 * time.Hour
	year = 365 * day  // approx
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

	payload, err := jwt.Decode("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxNTg2MzI4NDcxLCJpYXQiOjE1ODYzMjc4NzF9.lfRdglVTka7RY9_uVXP3r0WQbYvgRcoe_oYsTHQhaXGt5TZny5H56sod_ipeUKI59XAvATOhiYstCqumSZb8OTiwe2aVQVA6kTq03xgM9YScyzrVhzInf63bQtm-HtHHlMMT06APP8K7sHgOSD2ZtjtxSCAY9odnch6uKPPLaXFBZYQG2UBATGv9l3mkJB-iC7HYxar88tHX6MOvzZANpZZISvoCW0uz8jdu_c840FCPAxzoBdH0nRk-Gj6LV7zlbet9cXI6Q34JCQk4B6QirD-bPNVGa65zgFifKRkEV2zV0f0YLcDb4BjZH3F7WUa9MH0I7IHHYcVEPyOQaga7Fr5I7hxyeuKL-7MAxkWBgP1BpXHkQKQs1rDh4aZvMFbSwXcLFzB3z00aLcpcTNtlruxgTIPaGJk1pbN7rODp_fOFMmmzbEb52XUSmR5FFK8l7xYTwWn2XT83_r2BAKRpZsjGossiyVGAO3ahEJSr_TGlwYT7EM6sigIXwOeL9ivoIaaBzARHmoz4v9rOj3qCfdiRy3ruak3Ejz6YxmkfBDnIqMIyKjNnYZLfDOg5XtbPBY3d6ybbYx-yLfNIVnm6KlBk0Cy2fOYe3wwxJkwTMS869S3Iv3pbeUVMcX6Vd0z8R2DWXstKdLicywxjEjbt74hKvdlNZOjbAgygLJD797o")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, payload.Customer == "abc123", payload.Customer)
	assert.Assert(t, payload.RealUser == "xyz234", payload.RealUser)
	assert.Assert(t, payload.EffectiveUser == "xyz345", payload.EffectiveUser)
}
