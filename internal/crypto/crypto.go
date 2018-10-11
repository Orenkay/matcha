package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	tokenSecret = []byte("super-secret-key")
)

func EncryptPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	return string(hash), err
}

func ComparePassword(hash string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func CreateJWT(claims jwt.Claims) (string, *jwt.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", nil, err
	}
	return ss, token, nil
}

func DecodeJWT(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return tokenSecret, nil
	})
}

func RandToken(size int) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
