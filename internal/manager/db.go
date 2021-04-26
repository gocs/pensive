package manager

import "github.com/go-redis/redis"

type Manager struct {
	Cmdable redis.Cmdable
}

func NewManager(addr, password string) (*Manager, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
	})

	return &Manager{Cmdable: r}, r.Ping().Err()
}
