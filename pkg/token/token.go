package token

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

func Create(accessSecret string, userid uint64) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}
