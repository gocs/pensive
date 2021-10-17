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

// PostPublic is the post displayed through the site
type PostPublic struct {
	User           string
	Caption        string
	AttachmentURL  string
	AttachmentType string
	UpdatedAt      string
}

// updates list
// search post by user
// search user by post
