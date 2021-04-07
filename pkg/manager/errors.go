package manager

import "errors"

var (
	// ErrUsernameTaken common error on registration form when the username already existed
	ErrUsernameTaken = errors.New("username taken")

	// ErrUserNotFound common error on login form when the user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrTypeMismatch specific error for capturing type mismatch
	ErrTypeMismatch = errors.New("the type didn't match")

	// ErrNilClient gives error message when redis client variable is nil
	ErrNilClient = errors.New("client is nil")

	// ErrEmptyForm gives error message when form body is empty
	ErrEmptyForm = errors.New("form is empty")

	// ErrPerm gives error message when self tries to modify other accounts details
	ErrPerm = errors.New("permission to modify is denied")
)
