package manager

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gocs/pensive"
	sessions "github.com/gocs/pensive/internal/session"
	"github.com/gocs/pensive/pkg/timelayout"
	"github.com/redis/go-redis/v9"
)

// User the general people
type User struct {
	// make sure to pass redis pointer to c
	c  redis.Cmdable
	id int64
}

// ID UserID getter
// there should be no setter for id
func (u *User) ID() int64 { return u.id }

// Username Username getter
func (u *User) Username(ctx context.Context) (string, error) {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HGet(ctx, key, "username").Result()
}

// Password Hash getter
func (u *User) Password(ctx context.Context) ([]byte, error) {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HGet(ctx, key, "password").Bytes()
}

// Email getter
func (u *User) Email(ctx context.Context) (string, error) {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HGet(ctx, key, "email").Result()
}

// IsVerified getter
func (u *User) IsVerified(ctx context.Context) (bool, error) {
	key := fmt.Sprintf("user:%d", u.id)
	result, err := u.c.HGet(ctx, key, "is_verified").Result()
	if err != nil {
		return false, err
	}
	return result == "true", nil
}

// CreatedAt getter
// there should be no setter for created at
func (u *User) CreatedAt(ctx context.Context) (*time.Time, error) {
	key := fmt.Sprintf("user:%d", u.id)
	timeStr, err := u.c.HGet(ctx, key, "created_at").Result()
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(timelayout.UTCLayout, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, err
}

// UpdatedAt Hash getter
func (u *User) UpdatedAt(ctx context.Context) (*time.Time, error) {
	key := fmt.Sprintf("user:%d", u.id)
	timeStr, err := u.c.HGet(ctx, key, "updated_at").Result()
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(timelayout.UTCLayout, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, err
}

// SetUsername Username setter
func (u *User) SetUsername(ctx context.Context, value string) error {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HSet(ctx, key, "username", value).Err()
}

// SetPassword Password Hash setter
func (u *User) SetPassword(ctx context.Context, value []byte) error {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HSet(ctx, key, "password", value).Err()
}

// SetEmail Email setter
func (u *User) SetEmail(ctx context.Context, value string) error {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HSet(ctx, key, "email", value).Err()
}

// Verify Verify setter
// this converts bool to a redis friendly string bool
func (u *User) Verify(ctx context.Context, value bool) error {
	result := "false"
	if value {
		result = "true"
	}

	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HSet(ctx, key, "is_verified", result).Err()
}

// UpdateNow UpdatedAt setter
func (u *User) UpdateNow(ctx context.Context) error {
	now := time.Now().UTC().String()
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HSet(ctx, key, "updated_at", now).Err()
}

// AddUser this creates a new user entry
func AddUser(ctx context.Context, c redis.Cmdable, username string, password []byte, email string) (*User, error) {
	exists, err := c.HExists(ctx, "user:by-username", username).Result()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	// set created/updated at date
	now := time.Now().UTC().String()

	id, err := c.Incr(ctx, "user:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	pipe := c.Pipeline()
	pipe.HSet(ctx, key, "id", id)
	pipe.HSet(ctx, key, "username", username)
	pipe.HSet(ctx, key, "password", password)
	pipe.HSet(ctx, key, "email", email)
	pipe.HSet(ctx, key, "is_verified", "false")
	pipe.HSet(ctx, key, "created_at", now)
	pipe.HSet(ctx, key, "updated_at", now)
	pipe.HSet(ctx, "user:by-username", username, id)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &User{id: id, c: c}, nil
}

// RegisterUser register a valid user
func RegisterUser(ctx context.Context, c redis.Cmdable, username, password, email string) (*User, error) {
	hash, err := pensive.ValidatePassword(password)
	if err != nil {
		return nil, err
	}

	u, err := AddUser(ctx, c, username, hash, email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// AuthSelf checks if username has registered in this site
func AuthSelf(r *http.Request, s *sessions.Session, c redis.Cmdable, key string) (*User, error) {
	userID, err := s.GetInt64(r, key)
	if err != nil {
		return nil, err
	}

	u := &User{id: userID, c: c}
	un, err := u.Username(r.Context())
	if err != nil || un == "" {
		return nil, err
	}
	return u, nil
}

// GetUserByUserID gets user using a user id
func GetUserByUserID(c redis.Cmdable, userID int64) (*User, error) {
	if c == nil {
		return nil, redis.Nil
	}
	return &User{id: userID, c: c}, nil
}

// GetUserByName gets the user using the username
func GetUserByName(ctx context.Context, c redis.Cmdable, username string) (*User, error) {
	id, err := c.HGet(ctx, "user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return &User{id: id, c: c}, nil
}

func GetUser(ctx context.Context, c redis.Cmdable, user *User) (pensive.User, error) {
	username, err := user.Username(ctx)
	if err != nil {
		return pensive.User{}, err
	}

	password, err := user.Password(ctx)
	if err != nil {
		return pensive.User{}, err
	}

	email, err := user.Email(ctx)
	if err != nil {
		return pensive.User{}, err
	}

	createdAt, err := user.CreatedAt(ctx)
	if err != nil {
		return pensive.User{}, err
	}

	updatedAt, err := user.UpdatedAt(ctx)
	if err != nil {
		return pensive.User{}, err
	}

	u := pensive.User{
		ID:        user.ID(),
		Username:  username,
		Hash:      password,
		Email:     email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return u, nil
}

// AuthUser authenticates the user by its username and password
func AuthUser(ctx context.Context, c redis.Cmdable, username, password string) (*User, error) {
	user, err := GetUserByName(ctx, c, username)
	if err != nil {
		return nil, err
	}

	hash, err := user.Password(ctx)
	if err != nil {
		return nil, err
	}

	if err := pensive.Authenticate(hash, password); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser authenticates the user by its username and password and updates it with new ones
// userID is the current user requesting
func UpdateUser(ctx context.Context, c redis.Cmdable, userID int64, newUsername, oldPassword, newPassword, newEmail string) error {
	// get the account to be updated
	// if username is taken, make sure it is owned by the current user
	_, err := GetUserByName(ctx, c, newUsername)
	if err != nil {
		// skip vacant username
		if err != ErrUserNotFound {
			return err
		}
	}

	user := &User{id: userID, c: c}
	hash, err := user.Password(ctx)
	if err != nil {
		return err
	}
	// verify if the specified old password matched the current password
	if err := pensive.Authenticate(hash, oldPassword); err != nil {
		return err
	}

	oldUsername, err := user.Username(ctx)
	if err != nil {
		return err
	}

	// set updated at date
	now := time.Now().UTC().String()

	// update all details in the pipe
	key := fmt.Sprintf("user:%d", user.id)
	pipe := c.Pipeline()
	pipe.HSet(ctx, key, "id", user.id)
	if newUsername != "" {
		pipe.HSet(ctx, key, "username", newUsername)
	}
	if newPassword != "" {
		pipe.HSet(ctx, key, "password", newPassword)
	}
	if newEmail != "" {
		pipe.HSet(ctx, key, "email", newEmail)
	}
	pipe.HSet(ctx, key, "updated_at", now)
	pipe.HDel(ctx, "user:by-username", oldUsername)
	pipe.HSet(ctx, "user:by-username", newUsername, user.id)
	_, err = pipe.Exec(ctx)
	return err
}

// UpdateUsername authenticates the user by its username and password and updates it with new ones
// userID is the current user requesting
func UpdateUsername(ctx context.Context, c redis.Cmdable, userID int64, newUsername, oldPassword string) error {
	// get the account to be updated
	// if username is taken, make sure it is owned by the current user
	_, err := GetUserByName(ctx, c, newUsername)
	if err != nil {
		// skip vacant username
		if err != ErrUserNotFound {
			return err
		}
	}

	user := &User{id: userID, c: c}
	hash, err := user.Password(ctx)
	if err != nil {
		return err
	}
	// verify if the specified old password matched the current password
	if err := pensive.Authenticate(hash, oldPassword); err != nil {
		return err
	}

	oldUsername, err := user.Username(ctx)
	if err != nil {
		return err
	}

	// set updated at date
	now := time.Now().UTC().String()

	// update all details in the pipe
	key := fmt.Sprintf("user:%d", user.id)
	pipe := c.Pipeline()
	pipe.HSet(ctx, key, "id", user.id)
	if newUsername != "" {
		pipe.HSet(ctx, key, "username", newUsername)
	}
	pipe.HSet(ctx, key, "updated_at", now)
	pipe.HDel(ctx, "user:by-username", oldUsername)
	pipe.HSet(ctx, "user:by-username", newUsername, user.id)
	_, err = pipe.Exec(ctx)
	return err
}

// UpdatePassword authenticates the user by its username and password and updates it with new ones
// userID is the current user requesting
func UpdatePassword(ctx context.Context, c redis.Cmdable, userID int64, oldPassword, newPassword string) error {
	user := &User{id: userID, c: c}
	hash, err := user.Password(ctx)
	if err != nil {
		return err
	}
	// verify if the specified old password matched the current password
	if err := pensive.Authenticate(hash, oldPassword); err != nil {
		return err
	}

	// set updated at date
	now := time.Now().UTC().String()

	// update all details in the pipe
	key := fmt.Sprintf("user:%d", user.id)
	pipe := c.Pipeline()
	pipe.HSet(ctx, key, "id", user.id)
	if newPassword != "" {
		pipe.HSet(ctx, key, "password", newPassword)
	}
	pipe.HSet(ctx, key, "updated_at", now)
	_, err = pipe.Exec(ctx)
	return err
}

// UpdateEmail authenticates the user by its username and password and updates it with new ones
// userID is the current user requesting
func UpdateEmail(ctx context.Context, c redis.Cmdable, userID int64, newEmail, oldPassword string) error {
	user := &User{id: userID, c: c}
	hash, err := user.Password(ctx)
	if err != nil {
		return err
	}
	// verify if the specified old password matched the current password
	if err := pensive.Authenticate(hash, oldPassword); err != nil {
		return err
	}

	// set updated at date
	now := time.Now().UTC().String()

	// update all details in the pipe
	key := fmt.Sprintf("user:%d", user.id)
	pipe := c.Pipeline()
	pipe.HSet(ctx, key, "id", user.id)
	if newEmail != "" {
		pipe.HSet(ctx, key, "email", newEmail)
	}
	pipe.HSet(ctx, key, "updated_at", now)
	_, err = pipe.Exec(ctx)
	return err
}
