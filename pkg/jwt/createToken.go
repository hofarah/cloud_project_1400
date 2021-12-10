package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type Payload struct {
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("invalid token")

func (payload *Payload) Valid() error {
	if time.Now().Unix() > payload.Exp {
		return ErrExpiredToken
	}
	return nil
}

func CreateToken(username string) (string, error) {
	var err error
	payload := Payload{
		Username: username,
		Exp:      time.Now().Add(time.Hour * 48).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	key := os.Getenv("JWT_KEY")
	if key == "" {
		key = "HOMRSAYE"
	}
	token, err := t.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return token, nil
}
