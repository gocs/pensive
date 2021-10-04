package token

import (
	"fmt"
	"time"

	"github.com/gocs/errored"
	jwt "github.com/golang-jwt/jwt/v4"
)

// Create will generate a token to be verified
// signingKey should be auto-generated and should be kept secret
// userID will be used to verify users
func Create(signingKey string, userID string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// Verify 
// signingKey should be auto-generated and should be kept secret
// tokenStr is the string generated from the Create function
// userID is used to verify users
func Verify(signingKey, tokenStr string, userID string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errored.New("user is not authorized")
	}

	userIDClaims := fmt.Sprint(claims["user_id"])
	if userIDClaims != userID {
		return errored.New("claims doesn't match the expected user")
	}
	return nil
}
