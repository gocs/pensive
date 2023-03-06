package manager

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Manager struct {
	Cmdable redis.Cmdable
}

func NewManager(ctx context.Context, addr, password string) (*Manager, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,
	})

	if err := r.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Manager{Cmdable: r}, nil
}
