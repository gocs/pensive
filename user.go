package pensive

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  // required, nonzero
	Username string // required, nonzero
	Hash     []byte // required, nonzero
	Email    string // required, nonzero

	// do not manually fill-in the ff.
	CreatedAt *time.Time // required, nonzero
	UpdatedAt *time.Time // required, nonzero
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
