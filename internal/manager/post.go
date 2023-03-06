package manager

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gocs/pensive"
	"github.com/gocs/pensive/pkg/timelayout"
	"github.com/redis/go-redis/v9"
)

// Post user generated inputs
type Post struct {
	// make sure to pass redis pointer to c
	c  redis.Cmdable
	id int64
}

// Body getter
func (p *Post) Body(ctx context.Context) (string, error) {
	key := fmt.Sprintf("post:%d", p.id)
	return p.c.HGet(ctx, key, "body").Result()
}

// MediaID getter
func (p *Post) MediaID(ctx context.Context) (string, error) {
	key := fmt.Sprintf("post:%d", p.id)
	return p.c.HGet(ctx, key, "media_id").Result()
}

// User getter
func (p *Post) User(ctx context.Context) (*User, error) {
	key := fmt.Sprintf("post:%d", p.id)
	id, err := p.c.HGet(ctx, key, "user_id").Int64()
	if err != nil {
		return nil, err
	}

	return &User{id: id, c: p.c}, nil
}

// CreatedAt getter
func (p *Post) CreatedAt(ctx context.Context) (*time.Time, error) {
	key := fmt.Sprintf("post:%d", p.id)
	timeStr, err := p.c.HGet(ctx, key, "created_at").Result()
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
func (p *Post) UpdatedAt(ctx context.Context) (*time.Time, error) {
	key := fmt.Sprintf("post:%d", p.id)
	timeStr, err := p.c.HGet(ctx, key, "updated_at").Result()
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
func AddPost(ctx context.Context, c redis.Cmdable, p pensive.Post) (*Post, error) {
	id, err := c.Incr(ctx, "post:next-id").Result()
	if err != nil {
		return nil, err
	}

	// set created at date
	now := time.Now().UTC().String()

	key := fmt.Sprintf("post:%d", id)
	pipe := c.Pipeline()
	pipe.HSet(ctx, key, "id", id)
	pipe.HSet(ctx, key, "user_id", p.User.ID)
	pipe.HSet(ctx, key, "body", p.Body)
	pipe.HSet(ctx, key, "media_id", p.MediaID)
	pipe.HSet(ctx, key, "created_at", now)
	pipe.HSet(ctx, key, "updated_at", now)
	pipe.LPush(ctx, "posts", id)
	pipe.LPush(ctx, fmt.Sprintf("user:%d:posts", p.User.ID), id)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &Post{id: id, c: c}, nil
}

func queryPosts(ctx context.Context, c redis.Cmdable, key string) ([]*Post, error) {
	postIDs, err := c.LRange(ctx, key, 0, 10).Result()
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
func GetAllPosts(ctx context.Context, c redis.Cmdable) ([]*Post, error) {
	return queryPosts(ctx, c, "posts")
}

// GetPosts gets all updates created by the user
func GetPosts(ctx context.Context, c redis.Cmdable, userID int64) ([]*Post, error) {
	key := fmt.Sprintf("user:%d:posts", userID)
	return queryPosts(ctx, c, key)
}

// PostUpdate adds a new update; this differs from edit which actually changes
func PostUpdate(ctx context.Context, c redis.Cmdable, userID int64, body, mediaID string) error {
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

	_, err := AddPost(ctx, c, p)
	return err
}
