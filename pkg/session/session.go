package sessions

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	// ErrUserNotLoggedIn give error when session does not account you of logging in
	ErrUserNotLoggedIn = errors.New("you must login first")
)

// Session is used for storing session cookies
type Session struct {
	store       *sessions.CookieStore
	sessionName string
}

// New creates new session using a secret key
func New(secret, sessionName string) *Session {
	return &Session{
		store:       sessions.NewCookieStore([]byte(secret)),
		sessionName: sessionName,
	}
}

// Get is used to get a value of any type that is temporarily saved by a key, returns interface to be marshalled
func (s *Session) Get(r *http.Request, key string) (interface{}, error) {
	session, err := s.store.Get(r, s.sessionName)
	if err != nil {
		return "", err
	}

	return session.Values[key], nil
}

// GetInt64 is used to get value that is temporarily saved by a key, returns an int64 preferably the user id
func (s *Session) GetInt64(r *http.Request, key string) (int64, error) {
	session, err := s.store.Get(r, s.sessionName)
	if err != nil {
		return -1, err
	}

	u := session.Values[key]
	userID, ok := u.(int64)
	if !ok || u == nil {
		return -1, ErrUserNotLoggedIn
	}
	return userID, nil
}

func (s *Session) Set(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
	session, err := s.store.Get(r, s.sessionName)
	if err != nil {
		return err
	}
	session.Values[key] = value
	return session.Save(r, w)
}

func (s *Session) UnSet(w http.ResponseWriter, r *http.Request, key string) error {
	session, err := s.store.Get(r, s.sessionName)
	if err != nil {
		return err
	}
	delete(session.Values, key)
	return session.Save(r, w)
}
