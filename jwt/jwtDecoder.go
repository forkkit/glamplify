package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

type JwtDecoder struct {
	verifyKey *rsa.PublicKey
}

func NewJWTDecoderFromPath(pubKeyPath string) (JwtDecoder, error) {

	verifyBytes, _ := ioutil.ReadFile(pubKeyPath)
	return NewJWTDecoderFromBytes(verifyBytes)
}

func NewJWTDecoderFromBytes(verifyBytes []byte) (JwtDecoder, error) {

	verifyKey, err := jwtgo.ParseRSAPublicKeyFromPEM(verifyBytes)
	return JwtDecoder{
		verifyKey: verifyKey,
	}, err
}

func (jwt JwtDecoder) Decode(tokenString string) (Payload, error) {
	// sample token string taken from the New example
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	token, err := jwtgo.Parse(tokenString, func(token *jwtgo.Token) (interface{}, error) {
		return jwt.verifyKey, nil
	})

	data := Payload{}
	if err != nil {
		return data, err
	}

	if claims, ok := token.Claims.(jwtgo.MapClaims); ok && token.Valid {
		data.Customer, err = jwt.extractKey(claims, "accountId")
		if err != nil {
			return data, err
		}
		data.RealUser, err = jwt.extractKey(claims, "realUserId")
		if err != nil {
			return data, err
		}
		data.EffectiveUser, err = jwt.extractKey(claims, "effectiveUserId")
		if err != nil {
			return data, err
		}
		return data, nil
	}

	return data, errors.New("invalid claim token in jwt")
}

func (jwt JwtDecoder) extractKey(claims jwtgo.MapClaims, key string) (string, error) {

	val, ok := claims[key].(string)
	if !ok {
		return "", fmt.Errorf("missing %s in jwt token", key)
	}

	return val, nil
}
