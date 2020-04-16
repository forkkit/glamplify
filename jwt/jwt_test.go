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

	payload, err := jwt.Decode("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxNTg2OTk0Nzg0LCJpYXQiOjE1ODY5OTQxODR9.A7UQHRYpQ8b_jjV3Fy5mIrN9yamGZBc63ejuzuf5GNT3cbhkxAGhZ9qKRoqjQChqQY39Z_Rvj6pNZgJkayuKT5RhdvAYmOSzpCbq2sGHllikCIzJcedTNsId9qRgqw44ZITRJpe7xLJDjPaXsF2G3yHyb6hpqxxTFvowtTQIVdasedDL8OmGsbv7sKQUhukYgxKqeEK3DTG9DKqgkTpmcnNNxUxCuP-u69zXnB6548It0R_kn_lBM6yfIU08GmlJ5wNk8PRRsbAwtsq47_Du3-_8jRK-8kMXeZ2sp8CNwdawalgd2q0cOanl5ks4Iflg3PATdvcrLRl2j7po_Jt4mmX32QXf1-FuEPEemqLz96iaswemEIyekqbxSlpjsrd5EKlJAosAx3sh8_vOxA2o3cSL2xkG-I8oGI2YdOUQQNhLCarxVI0kAGzyfhOk9a2tygn0wj3I98ieo953qffolJyKbUKcrEXgHNrXKZM-JV2-2VLm00vcTIug4XXXr9r9hH4tZK52s8wftpe9VTmRqjOTAkmKTzn7SeiBKCiu2FncilU4K_yvfB4sTRT6YIcIV8PTkeYhIBURhixfmRMLRRbJHqePClg0uEY9ijeTz0-pys9qcoKnvRKFnGJKBlaVA8l6gXd_cmBBwhyzAF_brQT3KJAMQ2gKAKBJtMg2rMU")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, payload.Customer == "abc123", payload.Customer)
	assert.Assert(t, payload.RealUser == "xyz234", payload.RealUser)
	assert.Assert(t, payload.EffectiveUser == "xyz345", payload.EffectiveUser)
}
