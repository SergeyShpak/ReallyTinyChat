package jwt

import (
	"crypto/rand"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
)

func GenerateSecret() ([]byte, error) {
	len := 16
	secret := make([]byte, len)
	_, err := rand.Read(secret)
	if err != nil {

		return nil, errors.NewServerError(500, fmt.Sprintf("error when generating a JWT secret: ", err))
	}
	return secret, nil
}

func Verify(msg *types.Message, secret []byte) (string, error) {
	t, err := jwtgo.ParseWithClaims(msg.Token, &myCustomClaims{}, func(t *jwtgo.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", errors.NewServerError(401, fmt.Sprintf("error while parsing the JWT token: %v", err))
	}
	claims, ok := t.Claims.(*myCustomClaims)
	if !ok {
		return "", errors.NewServerError(401, "could not cast the JWT payload to the type with payload")
	}
	return claims.Payload, nil
}

type myCustomClaims struct {
	Payload string `json:"Payload"`
	jwtgo.StandardClaims
}
