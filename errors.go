package pensive

import "errors"

var (
	// ErrInvalidLogin common error on login form when the user does not match its login credentials
	ErrInvalidLogin = errors.New("invalid login")
)
