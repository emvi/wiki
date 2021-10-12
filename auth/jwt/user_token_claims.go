package jwt

import (
	"emviwiki/shared/constants"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"strings"
	"time"
)

type UserTokenClaims struct {
	jwt.StandardClaims

	UserId    hide.ID
	Language  string
	Scopes    []string
	IsSSOUser bool
}

func NewUserToken(claims *UserTokenClaims) (string, time.Time, error) {
	exp := time.Now().Add(constants.JwtSessionExp)
	claims.StandardClaims = getStandardClaims(exp)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(signKey)

	if err != nil {
		logbuch.Error("Error creating new user token", logbuch.Fields{"err": err})
		return "", time.Time{}, err
	}

	return tokenString, exp, nil
}

func getStandardClaims(exp time.Time) jwt.StandardClaims {
	now := time.Now()

	return jwt.StandardClaims{
		ExpiresAt: exp.Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
	}
}

func GetUserTokenClaims(token string) *UserTokenClaims {
	claims := new(UserTokenClaims)

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected token signing method: %v", token.Header["alg"])
		}

		return verifyKey, nil
	})

	if err != nil {
		if !strings.HasPrefix(err.Error(), tokenExpiredPrefix) {
			logbuch.Warn("Error parsing user JWT token", logbuch.Fields{"err": err, "token": token})
		}

		return nil
	}

	return claims
}
