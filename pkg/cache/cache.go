package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Define callback, when returning error
type Closure func(bytes []byte) error

const (
	cacheNil string = `redis: nil`
)

// AgentCache contract
type Cacher interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val interface{}, duration time.Duration) error
	Delete(ctx context.Context, key ...string) error
}

type cache struct {
	rds             redis.Cmdable
	retentionSecond time.Duration
}

// NewAgentCache creates new agent redis client
func NewCache(redis redis.Cmdable) Cacher {
	return &cache{
		rds: redis,
	}
}

func (c *cache) Set(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	cmd := c.rds.Set(ctx, key, val, exp)
	return cmd.Err()
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	cmd := c.rds.Get(ctx, key)
	b, e := cmd.Bytes()

	if e != nil {
		if e.Error() == cacheNil {
			return b, nil
		}
	}

	return b, e
}

func (c *cache) Delete(ctx context.Context, key ...string) error {
	cmd := c.rds.Del(ctx, key...)
	return cmd.Err()
}
