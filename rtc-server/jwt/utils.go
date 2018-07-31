package jwt

import (
	"crypto/rand"
	"fmt"
	"log"

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
	log.Println("Secret: ", secret)
	t, err := jwtgo.Parse(msg.Token, func(t *jwtgo.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", errors.NewServerError(401, fmt.Sprintf("error while parsing the JWT token: %v", err))
	}
	log.Println("Claims: ", t.Claims)
	return "", nil
}
