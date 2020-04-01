package jwt

import (
	"crypto/rsa"
	jwtgo "github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

type JwtEncoder struct {
	pemKey *rsa.PrivateKey
}

func NewJWTEncoderFromPath(pemKeyPath string) (JwtEncoder, error) {

	pemBytes, _ := ioutil.ReadFile(pemKeyPath)
	return NewJWTEncoderFromBytes(pemBytes)
}

func NewJWTEncoderFromBytes(pemBytes []byte) (JwtEncoder, error) {

	pemKey, err := jwtgo.ParseRSAPrivateKeyFromPEM(pemBytes)
	return JwtEncoder{
		pemKey: pemKey,
	}, err
}

func (jwt JwtEncoder) Encode(payload Payload) (string, error) {

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, jwtgo.MapClaims{
		"accountId":       payload.Customer,
		"realUserId":      payload.RealUser,
		"effectiveUserId": payload.EffectiveUser,
	})

	return token.SignedString(jwt.pemKey)
}
