package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(key, val string, timeout uint) error
	Get(key string) (string, error)
}

type RedisCache struct {
	rdb *redis.ClusterClient
}

func NewRedisCache(address string) *RedisCache {
	opts := &redis.ClusterOptions{
		Addrs: []string{address},
	}

	return &RedisCache{
		rdb: redis.NewClusterClient(opts),
	}
}

func (rc *RedisCache) Set(key, val string, timeout uint) error {
	return rc.rdb.Set(context.Background(), key, val, time.Duration(60)*time.Second).Err()
}

func (rc *RedisCache) Get(key string) (string, error) {
	return rc.rdb.Get(context.Background(), key).Result()
}
