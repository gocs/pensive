package pensive

import "time"

type Post struct {
	ID      int64 // required, nonzero
	User    User  // required, nonzero
	Body    string
	MediaID string

	// do not manually fill-in the ff.
	CreatedAt *time.Time // required, nonzero
	UpdatedAt *time.Time // required, nonzero
}

// updates list
// search post by user
// search user by post
