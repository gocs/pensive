package pensive

import "time"

type Post struct {
	ID        int64
	UserID    User
	Body      string
	CreatedAt *time.Time
	UpdatedAt  *time.Time
}

// updates list
// search post by user
// search user by post
