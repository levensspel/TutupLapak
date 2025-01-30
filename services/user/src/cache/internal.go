package cache

import (
	"context"
	"log"

	"github.com/bytedance/sonic"
)

func (d RedisCacheClient) set(ctx context.Context, key string, value string) error {
	_, err := d.client.Set(ctx, key, string(value), DefaultCacheTtl).Result()
	if err != nil {
		return err
	}
	return nil
}

func (d RedisCacheClient) setAsMap(ctx context.Context, key string, value map[string]string) error {
	data, err := sonic.Marshal(value)
	if err != nil {
		panic(err)
	}

	_, err = d.client.Set(ctx, key, string(data), DefaultCacheTtl).Result()
	if err != nil {
		return err
	}
	return nil
}

func (d RedisCacheClient) get(ctx context.Context, key string) (string, bool) {
	result, err := d.client.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}

	return result, true
}

func (d RedisCacheClient) getAsMap(ctx context.Context, key string) (map[string]string, bool) {
	val, err := d.client.Get(ctx, key).Result()

	if err != nil {
		return nil, false
	}

	var result map[string]string
	err = sonic.Unmarshal([]byte(val), &result)
	if err != nil {
		log.Printf("failed to unmarshal cache value: %v", val)
		panic(err)
	}

	return result, true
}

func (d RedisCacheClient) getAsMapArray(ctx context.Context, key string) ([]map[string]string, bool) {
	val, err := d.client.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	var result []map[string]string
	err = sonic.Unmarshal([]byte(val), &result)
	if err != nil {
		panic(err)
	}

	return result, true
}
