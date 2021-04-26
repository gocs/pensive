package manager

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gocs/pensive"
	"github.com/gocs/pensive/pkg/timelayout"
)

// Post user generated inputs
type Post struct {
	// make sure to pass redis pointer to c
	c  redis.Cmdable
	id int64
}

// Body getter
func (p *Post) Body() (string, error) {
	key := fmt.Sprintf("post:%d", p.id)
	return p.c.HGet(key, "body").Result()
}

// MediaID getter
func (p *Post) MediaID() (string, error) {
	key := fmt.Sprintf("post:%d", p.id)
	return p.c.HGet(key, "media_id").Result()
}

// User getter
func (p *Post) User() (*User, error) {
	key := fmt.Sprintf("post:%d", p.id)
	id, err := p.c.HGet(key, "user_id").Int64()
	if err != nil {
		return nil, err
	}

	return &User{id: id, c: p.c}, nil
}

// CreatedAt getter
func (p *Post) CreatedAt() (*time.Time, error) {
	key := fmt.Sprintf("post:%d", p.id)
	timeStr, err := p.c.HGet(key, "created_at").Result()
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(timelayout.UTCLayout, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, err
}

// UpdatedAt getter
func (p *Post) UpdatedAt() (*time.Time, error) {
	key := fmt.Sprintf("post:%d", p.id)
	timeStr, err := p.c.HGet(key, "updated_at").Result()
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(timelayout.UTCLayout, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, err
}

// AddPost creates a new post, saves it to the database, and returns the newly created question
func AddPost(c redis.Cmdable, p pensive.Post) (*Post, error) {
	id, err := c.Incr("post:next-id").Result()
	if err != nil {
		return nil, err
	}

	// set created at date
	now := time.Now().UTC().String()

	key := fmt.Sprintf("post:%d", id)
	pipe := c.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "user_id", p.User.ID)
	pipe.HSet(key, "body", p.Body)
	pipe.HSet(key, "media_id", p.MediaID)
	pipe.HSet(key, "created_at", now)
	pipe.HSet(key, "updated_at", now)
	pipe.LPush("posts", id)
	pipe.LPush(fmt.Sprintf("user:%d:posts", p.ID), id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}

	return &Post{id: id, c: c}, nil
}

func queryPosts(c redis.Cmdable, key string) ([]*Post, error) {
	postIDs, err := c.LRange(key, 0, 10).Result()
	if err != nil {
		return nil, err
	}

	posts := make([]*Post, len(postIDs))
	for i, val := range postIDs {
		id, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		posts[i] = &Post{id: id, c: c}
	}
	return posts, nil
}

// GetAllPosts All Updates getter
func GetAllPosts(c redis.Cmdable) ([]*Post, error) {
	return queryPosts(c, "posts")
}

// GetPosts gets all updates created by the user
func GetPosts(c redis.Cmdable, userID int64) ([]*Post, error) {
	key := fmt.Sprintf("user:%d:posts", userID)
	return queryPosts(c, key)
}

// PostUpdate adds a new update; this differs from edit which actually changes
func PostUpdate(c redis.Cmdable, userID int64, body, mediaID string) error {
	if body == "" {
		if mediaID == "" {
			return ErrEmptyForm
		}
	}

	p := pensive.Post{
		User:    pensive.User{ID: userID},
		Body:    body,
		MediaID: mediaID,
	}

	_, err := AddPost(c, p)
	return err
}
