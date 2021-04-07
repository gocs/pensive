package pensive

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Username  string
	Hash      []byte
	Email     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func ValidatePassword(hash string) ([]byte, error) {
	cost := bcrypt.DefaultCost
	return bcrypt.GenerateFromPassword([]byte(hash), cost)
}

func Authenticate(hash []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidLogin
	}
	return err
}
