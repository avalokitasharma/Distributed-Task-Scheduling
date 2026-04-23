package cache

import (
	"context"
	"sync"

	"github.com/avalokitasharma/job-scheduler/pkg/config"
	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

func Get() *redis.Client {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr: config.App.RedisAddr,
		})
		client.Ping(context.Background())
	})
	return client
}
