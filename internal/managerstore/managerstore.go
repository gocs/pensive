package managerstore

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/gocs/pensive"
	"github.com/gocs/pensive/internal/manager"
	"github.com/gocs/pensive/pkg/file"
	"github.com/gocs/pensive/pkg/objectstore"
)

func ListPost(ctx context.Context, objs *objectstore.ObjectStore, c redis.Cmdable) ([]pensive.PostPublic, error) {
	posts, err := manager.GetAllPosts(c)
	if err != nil {
		return nil, err
	}

	ps := []pensive.PostPublic{}
	for _, post := range posts {
		user, err := post.User()
		if err != nil {
			return nil, err
		}
		username, err := user.Username()
		if err != nil {
			return nil, err
		}
		body, err := post.Body()
		if err != nil {
			return nil, err
		}
		filename, err := post.MediaID()
		if err != nil {
			return nil, err
		}
		updatedAt, err := post.UpdatedAt()
		if err != nil {
			return nil, err
		}

		bName := fmt.Sprintf("user%d", user.ID())

		attachmentURL := ""
		if filename != "" {
			opts := objectstore.PresignedGetObjectOptions{}
			url, err := objs.GetPresignedURLObject(ctx, bName, filename, opts)
			if err != nil {
				return nil, err
			}

			attachmentURL = url.String()
		}

		ps = append(ps, pensive.PostPublic{
			User:           username,
			Caption:        body,
			AttachmentURL:  attachmentURL,
			AttachmentType: file.GetMediaType(filename),
			UpdatedAt:      updatedAt.Format(time.RFC822Z),
		})
	}
	return ps, nil
}

func ListPostByUserID(ctx context.Context, objs *objectstore.ObjectStore, c redis.Cmdable, userID int64) ([]pensive.PostPublic, error) {
	posts, err := manager.GetPosts(c, userID)
	if err != nil {
		return nil, err
	}

	ps := []pensive.PostPublic{}
	for _, post := range posts {
		user, err := post.User()
		if err != nil {
			return nil, err
		}
		username, err := user.Username()
		if err != nil {
			return nil, err
		}
		body, err := post.Body()
		if err != nil {
			return nil, err
		}
		filename, err := post.MediaID()
		if err != nil {
			return nil, err
		}
		updatedAt, err := post.UpdatedAt()
		if err != nil {
			return nil, err
		}

		bName := fmt.Sprintf("user%d", user.ID())

		attachmentURL := ""
		if filename != "" {
			opts := objectstore.PresignedGetObjectOptions{}
			url, err := objs.GetPresignedURLObject(ctx, bName, filename, opts)
			if err != nil {
				return nil, err
			}

			attachmentURL = url.String()
		}

		ps = append(ps, pensive.PostPublic{
			User:           username,
			Caption:        body,
			AttachmentURL:  attachmentURL,
			AttachmentType: file.GetMediaType(filename),
			UpdatedAt:      updatedAt.Format(time.RFC822Z),
		})
	}
	return ps, nil
}
