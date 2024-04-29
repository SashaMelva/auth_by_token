package pkg

import (
	"math/rand"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAccssesToken(userId string, jwtKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Unix(),
		Subject:   userId,
	})

	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken() ([]byte, error) {
	var b []byte
	timeSource := rand.NewSource(time.Now().Unix())
	reader := rand.New(timeSource)
	_, err := reader.Read(b)

	if err != nil {
		return []byte{}, err
	}

	return bcrypt.GenerateFromPassword([]byte(b), bcrypt.DefaultCost)
}
