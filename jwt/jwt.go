package jwt

import (
	"errors"
	"net/http"
	"strings"
)

type Payload struct {
	Customer      string // uuid
	RealUser      string // uuid
	EffectiveUser string // uid
}

type DecodeJwtToken interface {
	Decode(tokenString string) (Payload, error)
}

func PayloadFromRequest(r *http.Request, jwtDecoder DecodeJwtToken) (Payload, error) {
	auth := r.Header.Get("Authorization") // "Authorization: Bearer xxxxx.yyyyy.zzzzz"
	if len(auth) == 0 {
		return Payload{}, errors.New("missing authorization header")
	}

	splitToken := strings.Split(auth, "Bearer")
	if len(splitToken) < 2 {
		return Payload{}, errors.New("missing 'Bearer' token in authorization header")
	}

	token := strings.TrimSpace(splitToken[1])
	return jwtDecoder.Decode(token)
}

