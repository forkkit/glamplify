package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
)

type Decoder struct {
	verifyKey *rsa.PublicKey
}

func NewDecoder() (Decoder, error) {

	pubKey := os.Getenv("AUTH_PUBLIC_KEY")
	return NewDecoderFromBytes([]byte(pubKey))
}

func NewDecoderFromPath(pubKeyPath string) (Decoder, error) {

	verifyBytes, _ := ioutil.ReadFile(pubKeyPath)
	return NewDecoderFromBytes(verifyBytes)
}

func NewDecoderFromBytes(verifyBytes []byte) (Decoder, error) {

	verifyKey, err := jwtgo.ParseRSAPublicKeyFromPEM(verifyBytes)
	return Decoder{
		verifyKey: verifyKey,
	}, err
}

func (jwt Decoder) Decode(tokenString string) (Payload, error) {
	// sample token string in the form "header.payload.signature"
	//eg. "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	data := Payload{}

	token, err := jwtgo.Parse(tokenString, func(token *jwtgo.Token) (interface{}, error) {
		return jwt.verifyKey, nil
	})
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

func (jwt Decoder) extractKey(claims jwtgo.MapClaims, key string) (string, error) {

	val, ok := claims[key].(string)
	if !ok {
		return "", fmt.Errorf("missing %s in jwt token", key)
	}

	return val, nil
}
