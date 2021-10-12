package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emvi/logbuch"
	"strings"
	"time"
)

const (
	JwtClientExp       = time.Minute * 10
	tokenExpiredPrefix = "token is expired"
)

type ClientTokenClaims struct {
	jwt.StandardClaims

	ClientId string
	Trusted  bool
	Scopes   []string
}

func NewClientToken(claims *ClientTokenClaims) (string, time.Time, error) {
	exp := time.Now().Add(JwtClientExp)
	claims.StandardClaims = getStandardClaims(exp)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(signKey)

	if err != nil {
		logbuch.Error("Error creating new client token", logbuch.Fields{"err": err})
		return "", time.Time{}, err
	}

	return tokenString, exp, nil
}

func GetClientTokenClaims(token string) *ClientTokenClaims {
	claims := new(ClientTokenClaims)

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected token signing method: %v", token.Header["alg"])
		}

		return verifyKey, nil
	})

	if err != nil {
		if !strings.HasPrefix(err.Error(), tokenExpiredPrefix) {
			logbuch.Debug("Error parsing client JWT token", logbuch.Fields{"err": err, "token": token})
		}

		return nil
	}

	return claims
}
