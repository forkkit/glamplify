package jwt

import (
	"context"
	"crypto/rsa"
	"github.com/cultureamp/glamplify/log"
	jwtgo "github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

type JwtEncoder struct {
	pemKey *rsa.PrivateKey
	logger *log.Logger
}

func NewJWTEncoderFromPath(ctx context.Context, pemKeyPath string) JwtEncoder {

	logger := log.New(ctx)

	pemBytes, err := ioutil.ReadFile(pemKeyPath)
	if err != nil {
		logger.Error(err, log.Fields{
			"pem_key_path": pemKeyPath,
		})
	}

	return NewJWTEncoderFromBytes(ctx, pemBytes)
}

func NewJWTEncoderFromBytes(ctx context.Context, pemBytes []byte) JwtEncoder {
	logger := log.New(ctx)

	pemKey, err := jwtgo.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		logger.Error(err)
	}

	return JwtEncoder{
		pemKey: pemKey,
		logger:    logger,
	}
}

func (jwt JwtEncoder) Encode(payload Payload) (string, error) {

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, jwtgo.MapClaims{
		"accountId":       payload.Customer,
		"realUserId":      payload.RealUser,
		"effectiveUserId": payload.EffectiveUser,
	})

	tokenString, err := token.SignedString(jwt.pemKey)
	if err != nil {
		jwt.logger.Error(err)
		return "", err
	}

	return tokenString, nil
}
