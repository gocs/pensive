package managerstore

import (
	"context"
	"fmt"
	"time"

	"github.com/gocs/pensive"
	"github.com/gocs/pensive/internal/manager"
	"github.com/gocs/pensive/pkg/file"
	"github.com/gocs/pensive/pkg/objectstore"
	"github.com/redis/go-redis/v9"
)

func ListPost(ctx context.Context, objs *objectstore.ObjectStore, c redis.Cmdable) ([]pensive.PostPublic, error) {
	posts, err := manager.GetAllPosts(ctx, c)
	if err != nil {
		return nil, err
	}

	ps := []pensive.PostPublic{}
	for _, post := range posts {
		user, err := post.User(ctx)
		if err != nil {
			return nil, err
		}
		username, err := user.Username(ctx)
		if err != nil {
			return nil, err
		}
		body, err := post.Body(ctx)
		if err != nil {
			return nil, err
		}
		filename, err := post.MediaID(ctx)
		if err != nil {
			return nil, err
		}
		updatedAt, err := post.UpdatedAt(ctx)
		if err != nil {
			return nil, err
		}

		attachmentURL := ""
		if filename != "" {
			attachmentURL = fmt.Sprintf("/@%s/%s", username, filename)
		}

		ps = append(ps, pensive.PostPublic{
			User:           username,
			Caption:        body,
			AttachmentURL:  attachmentURL,
			AttachmentType: file.GetMediaType(filename),
			UpdatedAt:      updatedAt.Format(time.RFC822),
		})
	}
	return ps, nil
}

func ListPostByUserID(ctx context.Context, objs *objectstore.ObjectStore, c redis.Cmdable, userID int64) ([]pensive.PostPublic, error) {
	posts, err := manager.GetPosts(ctx, c, userID)
	if err != nil {
		return nil, err
	}

	ps := []pensive.PostPublic{}
	for _, post := range posts {
		user, err := post.User(ctx)
		if err != nil {
			return nil, err
		}
		username, err := user.Username(ctx)
		if err != nil {
			return nil, err
		}
		body, err := post.Body(ctx)
		if err != nil {
			return nil, err
		}
		filename, err := post.MediaID(ctx)
		if err != nil {
			return nil, err
		}
		updatedAt, err := post.UpdatedAt(ctx)
		if err != nil {
			return nil, err
		}

		attachmentURL := "#"
		if filename != "" {
			attachmentURL = fmt.Sprintf("/@%s/%s", username, filename)
		}

		ps = append(ps, pensive.PostPublic{
			User:           username,
			Caption:        body,
			AttachmentURL:  attachmentURL,
			AttachmentType: file.GetMediaType(filename),
			UpdatedAt:      updatedAt.Format(time.RFC822),
		})
	}
	return ps, nil
}
