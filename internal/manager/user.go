package manager

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gocs/pensive"
	"github.com/gocs/pensive/pkg/timelayout"
	sessions "github.com/gocs/pensive/internal/session"
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
func (u *User) Username() (string, error) {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HGet(key, "username").Result()
}

// Password Hash getter
func (u *User) Password() ([]byte, error) {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HGet(key, "password").Bytes()
}

// Email Hash getter
func (u *User) Email() (string, error) {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HGet(key, "email").Result()
}

// CreatedAt Hash getter
// there should be no setter for created at
func (u *User) CreatedAt() (*time.Time, error) {
	key := fmt.Sprintf("user:%d", u.id)
	timeStr, err := u.c.HGet(key, "created_at").Result()
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
func (u *User) UpdatedAt() (*time.Time, error) {
	key := fmt.Sprintf("user:%d", u.id)
	timeStr, err := u.c.HGet(key, "updated_at").Result()
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(timelayout.UTCLayout, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, err
}

// Username Username setter
func (u *User) SetUsername(value string) error {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HSet(key, "username", value).Err()
}

// Password Hash setter
func (u *User) SetPassword(value []byte) error {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HSet(key, "password", value).Err()
}

// Email Hash setter
func (u *User) SetEmail(value string) error {
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HGet(key, "email").Err()
}

// UpdatedAt Hash setter
func (u *User) UpdateNow() error {
	now := time.Now().UTC().String()
	key := fmt.Sprintf("user:%d", u.id)
	return u.c.HSet(key, "updated_at", now).Err()
}

// AddUser this creates a new user entry
func AddUser(c redis.Cmdable, username string, password []byte, email string) (*User, error) {
	exists, err := c.HExists("user:by-username", username).Result()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	// set created/updated at date
	now := time.Now().UTC().String()

	id, err := c.Incr("user:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	pipe := c.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "username", username)
	pipe.HSet(key, "password", password)
	pipe.HSet(key, "email", email)
	pipe.HSet(key, "created_at", now)
	pipe.HSet(key, "updated_at", now)
	pipe.HSet("user:by-username", username, id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &User{id: id, c: c}, nil
}

// RegisterUser register a valid user
func RegisterUser(c redis.Cmdable, username, password, email string) error {
	hash, err := pensive.ValidatePassword(password)
	if err != nil {
		return err
	}

	_, err = AddUser(c, username, hash, email)
	return err
}

// AuthSelf checks if username has registered in this site
func AuthSelf(r *http.Request, s *sessions.Session, c redis.Cmdable, key string) (*User, error) {
	userID, err := s.GetInt64(r, key)
	if err != nil {
		return nil, err
	}

	u := &User{id: userID, c: c}
	un, err := u.Username()
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
func GetUserByName(c redis.Cmdable, username string) (*User, error) {
	id, err := c.HGet("user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return &User{id: id, c: c}, nil
}

func GetUser(c redis.Cmdable, user *User) (pensive.User, error) {
	username, err := user.Username()
	if err != nil {
		return pensive.User{}, err
	}

	password, err := user.Password()
	if err != nil {
		return pensive.User{}, err
	}

	email, err := user.Email()
	if err != nil {
		return pensive.User{}, err
	}

	createdAt, err := user.CreatedAt()
	if err != nil {
		return pensive.User{}, err
	}

	updatedAt, err := user.UpdatedAt()
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
func AuthUser(c redis.Cmdable, username, password string) (*User, error) {
	user, err := GetUserByName(c, username)
	if err != nil {
		return nil, err
	}

	hash, err := user.Password()
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
func UpdateUser(c redis.Cmdable, userID int64, newUsername, oldPassword, newPassword, newEmail string) error {
	// get the account to be updated
	// if username is taken, make sure it is owned by the current user
	_, err := GetUserByName(c, newUsername)
	if err != nil {
		// skip vacant username
		if err != ErrUserNotFound {
			return err
		}
	}

	user := &User{id: userID, c: c}
	hash, err := user.Password()
	if err != nil {
		return err
	}
	// verify if the specified old password matched the current password
	if err := pensive.Authenticate(hash, oldPassword); err != nil {
		return err
	}

	oldUsername, err := user.Username()
	if err != nil {
		return err
	}

	// set updated at date
	now := time.Now().UTC().String()

	// update all details in the pipe
	key := fmt.Sprintf("user:%d", user.id)
	pipe := c.Pipeline()
	pipe.HSet(key, "id", user.id)
	if newUsername != "" {
		pipe.HSet(key, "username", newUsername)
	}
	if newPassword != "" {
		pipe.HSet(key, "password", newPassword)
	}
	if newEmail != "" {
		pipe.HSet(key, "email", newEmail)
	}
	pipe.HSet(key, "updated_at", now)
	pipe.HDel("user:by-username", oldUsername)
	pipe.HSet("user:by-username", newUsername, user.id)
	_, err = pipe.Exec()
	return err
}

// UpdateUsername authenticates the user by its username and password and updates it with new ones
// userID is the current user requesting
func UpdateUsername(c redis.Cmdable, userID int64, newUsername, oldPassword string) error {
	// get the account to be updated
	// if username is taken, make sure it is owned by the current user
	_, err := GetUserByName(c, newUsername)
	if err != nil {
		// skip vacant username
		if err != ErrUserNotFound {
			return err
		}
	}

	user := &User{id: userID, c: c}
	hash, err := user.Password()
	if err != nil {
		return err
	}
	// verify if the specified old password matched the current password
	if err := pensive.Authenticate(hash, oldPassword); err != nil {
		return err
	}

	oldUsername, err := user.Username()
	if err != nil {
		return err
	}

	// set updated at date
	now := time.Now().UTC().String()

	// update all details in the pipe
	key := fmt.Sprintf("user:%d", user.id)
	pipe := c.Pipeline()
	pipe.HSet(key, "id", user.id)
	if newUsername != "" {
		pipe.HSet(key, "username", newUsername)
	}
	pipe.HSet(key, "updated_at", now)
	pipe.HDel("user:by-username", oldUsername)
	pipe.HSet("user:by-username", newUsername, user.id)
	_, err = pipe.Exec()
	return err
}

// UpdatePassword authenticates the user by its username and password and updates it with new ones
// userID is the current user requesting
func UpdatePassword(c redis.Cmdable, userID int64, oldPassword, newPassword string) error {
	user := &User{id: userID, c: c}
	hash, err := user.Password()
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
	pipe.HSet(key, "id", user.id)
	if newPassword != "" {
		pipe.HSet(key, "password", newPassword)
	}
	pipe.HSet(key, "updated_at", now)
	_, err = pipe.Exec()
	return err
}

// UpdateEmail authenticates the user by its username and password and updates it with new ones
// userID is the current user requesting
func UpdateEmail(c redis.Cmdable, userID int64, newEmail, oldPassword string) error {
	user := &User{id: userID, c: c}
	hash, err := user.Password()
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
	pipe.HSet(key, "id", user.id)
	if newEmail != "" {
		pipe.HSet(key, "email", newEmail)
	}
	pipe.HSet(key, "updated_at", now)
	_, err = pipe.Exec()
	return err
}
