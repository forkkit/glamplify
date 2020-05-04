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

func PayloadFromRequest(r *http.Request) (Payload, error) {

	token := r.Header.Get("Authorization") // "Authorization: Bearer xxxxx.yyyyy.zzzzz"
	if len(token) == 0 {
		return Payload{}, errors.New("Missing authorization header")
	}

	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) < 2 {
		return Payload{}, errors.New("Missing 'Bearer' token in authorization header")
	}

	jwt, err := NewDecoder()
	if err != nil {
		return Payload{}, err
	}

	return jwt.Decode(splitToken[1])
}
