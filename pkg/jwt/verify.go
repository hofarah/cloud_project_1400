package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
)

func Verify(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		key := os.Getenv("JWT_KEY")
		if key == "" {
			key = "HOMRSAYE"
		}
		return []byte(key), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok && !errors.Is(validationErr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		if !errors.Is(validationErr.Inner, ErrExpiredToken) {
			return nil, ErrInvalidToken
		}
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
