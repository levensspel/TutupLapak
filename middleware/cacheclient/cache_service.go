package serviceCache

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/samber/do/v2"
	"log"
	"time"
)

type CacheService interface {
	Set(ctx context.Context, key string, value string) error
	SetAsMap(ctx context.Context, key string, value map[string]string) error
	Get(ctx context.Context, key string) (string, bool)
	GetAsMap(ctx context.Context, key string) (map[string]string, bool)
	GetAsMapArray(ctx context.Context, key string) ([]map[string]string, bool)
}

const (
	DefaultCacheTtl = 5 * time.Minute
)

type CacheServiceImpl struct {
	cacheClient CacheClient
}

func (c CacheServiceImpl) Set(ctx context.Context, key string, value string) error {
	cost := int64(len(key) + len(value))
	return c.cacheClient.Set(ctx, key, value, cost, DefaultCacheTtl)
}

func (c CacheServiceImpl) SetAsMap(ctx context.Context, key string, value map[string]string) error {
	data, err := sonic.Marshal(value)
	if err != nil {
		panic(err)
	}

	cost := int64(len(key) + len(data))
	return c.cacheClient.Set(ctx, key, string(data), cost, DefaultCacheTtl)
}

func (c CacheServiceImpl) Get(ctx context.Context, key string) (string, bool) {
	return c.cacheClient.Get(ctx, key)
}

func (c CacheServiceImpl) GetAsMap(ctx context.Context, key string) (map[string]string, bool) {
	val, found := c.cacheClient.Get(ctx, key)

	if !found {
		return nil, false
	}

	var result map[string]string
	err := sonic.Unmarshal([]byte(val), &result)
	if err != nil {
		log.Printf("failed to unmarshal cache value: %v", val)
		panic(err)
	}

	return result, true
}

func (c CacheServiceImpl) GetAsMapArray(ctx context.Context, key string) ([]map[string]string, bool) {
	val, found := c.cacheClient.Get(ctx, key)
	if !found {
		return nil, false
	}

	var result []map[string]string
	err := sonic.Unmarshal([]byte(val), &result)
	if err != nil {
		panic(err)
	}

	return result, true
}

func NewCacheService(cacheClient CacheClient) CacheService {
	return CacheServiceImpl{
		cacheClient: cacheClient,
	}
}

func NewCacheServiceInject(i do.Injector) (CacheService, error) {
	return NewCacheService(
		do.MustInvoke[CacheClient](i),
	), nil
}
