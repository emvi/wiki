package jwt

import (
	"crypto/rsa"
	"emviwiki/shared/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/emvi/logbuch"
	"io/ioutil"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func LoadRSAKeys() {
	publicKeyFile := config.Get().JWT.PublicKey
	privateKeyFile := config.Get().JWT.PrivateKey
	publicKey, err := ioutil.ReadFile(publicKeyFile)

	if err != nil {
		logbuch.Fatal("Token public key file not found", logbuch.Fields{"err": err, "file": publicKeyFile})
	}

	privateKey, err := ioutil.ReadFile(privateKeyFile)

	if err != nil {
		logbuch.Fatal("Token private key file not found", logbuch.Fields{"err": err, "file": privateKeyFile})
	}

	verify, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)

	if err != nil {
		logbuch.Fatal("Error parsing token public key", logbuch.Fields{"err": err})
	}

	sign, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	if err != nil {
		logbuch.Fatal("Error parsing token private key", logbuch.Fields{"err": err})
	}

	verifyKey = verify
	signKey = sign
}
