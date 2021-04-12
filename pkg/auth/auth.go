package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
}

type tokenClaims struct {
	jwt.StandardClaims
	username string
}

const (
	solt       = "dsfjdewkd3289djsoidwe"
	signingKey = "fsfdsfewr435re3wrfd43R$#443Swww"
	tokenTTL   = 12 * time.Hour
)

func (t *Auth) CreateToken(user string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix()},
		user,
	})
	return token.SignedString([]byte(signingKey))
}

func (t *Auth) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err

	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.username, nil

}
