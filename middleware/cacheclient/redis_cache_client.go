package serviceCache

import (
	"context"
	"fmt"
	"github.com/samber/do/v2"
	"log"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisCacheClientOnce     sync.Once
	redisCacheClientInstance *RedisCacheClient
)

type RedisCacheClient struct {
	client *redis.Client
}

func (d RedisCacheClient) Get(ctx context.Context, key string) (string, bool) {
	result, err := d.client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("failed to get cache value: %v", err)
		return "", false
	}
	return result, true
}

func (d RedisCacheClient) Set(
	ctx context.Context,
	key string,
	value string,
	cost int64,
	ttl time.Duration,
) error {
	_, err := d.client.Set(ctx, key, value, ttl).Result()
	if err != nil {
		log.Printf("failed to set cache value: %v", err)
		return err
	}
	return nil
}

func NewRedisCacheClient() CacheClient {
	redisCacheClientOnce.Do(func() {
		redisCacheClientInstance = &RedisCacheClient{
			client: redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf(
					"%s:%s",
					os.Getenv("REDIS_HOST"),
					os.Getenv("REDIS_PORT"),
				),
				Password: "", // for the moment, password is not required
				DB:       0,  // Default database number
			}),
		}
	})

	return redisCacheClientInstance
}

func NewRedisCacheClientInject(i do.Injector) (CacheClient, error) {
	return NewRedisCacheClient(), nil
}
