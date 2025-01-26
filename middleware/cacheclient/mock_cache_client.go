package serviceCache

import (
	"context"
	"fmt"
	"github.com/samber/do/v2"
	"time"
)

type MockCacheClient struct{}

func (m MockCacheClient) Get(ctx context.Context, key string) (string, bool) {
	if key == "mock_not_found" {
		return "", false
	}
	if key == "mock_as_map" {
		return `{"key1": "value1", "key2": "value2"}`, true
	}
	return fmt.Sprintf("value_of_%s", key), false
}

func (m MockCacheClient) Set(
	ctx context.Context,
	key string,
	value string,
	cost int64,
	ttl time.Duration,
) error {
	if key == "mock_error" {
		return fmt.Errorf("mock error set cache")
	}
	return nil
}

func NewMockCacheClient() CacheClient {
	return &MockCacheClient{}
}

func NewMockCacheClientInject(i do.Injector) (CacheClient, error) {
	return NewMockCacheClient(), nil
}
