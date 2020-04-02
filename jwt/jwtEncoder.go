package jwt

import (
	"crypto/rsa"
	jwtgo "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
	"time"
)

type Encoder struct {
	pemKey *rsa.PrivateKey
}

// Claims contains the claims to be used to sign JWT's returned by Identity API
type claims struct {
	AccountID       string `json:"accountId"`
	EffectiveUserID string `json:"effectiveUserId"`
	RealUserID      string `json:"realUserId"`
	jwtgo.StandardClaims
}

func NewJWTEncoder() (Encoder, error) {

	priKey := os.Getenv("AUTH_PRIVATE_KEY")
	return NewJWTEncoderFromBytes([]byte(priKey))
}

func NewJWTEncoderFromPath(pemKeyPath string) (Encoder, error) {

	pemBytes, _ := ioutil.ReadFile(pemKeyPath)
	return NewJWTEncoderFromBytes(pemBytes)
}

func NewJWTEncoderFromBytes(pemBytes []byte) (Encoder, error) {

	pemKey, err := jwtgo.ParseRSAPrivateKeyFromPEM(pemBytes)
	return Encoder{
		pemKey: pemKey,
	}, err
}

func (jwt Encoder) Encode(payload Payload) (string, error) {

	claims := jwt.claims(payload)
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, claims)
	return token.SignedString(jwt.pemKey)
}

func (jwt Encoder) claims(payload Payload) claims {
	now := time.Now()
	return claims{
		AccountID:       payload.Customer,
		EffectiveUserID: payload.EffectiveUser,
		RealUserID:      payload.RealUser,
		StandardClaims: jwtgo.StandardClaims{
			IssuedAt: now.Unix(),
			// Were a little loose on the expiry for now, to avoid possible
			// problems with clock skew, slow requests, background jobs (?) etc.
			ExpiresAt: now.Add(10 * time.Minute).Unix(),
		},
	}
}