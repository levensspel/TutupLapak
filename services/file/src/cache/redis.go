package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/TimDebug/TutupLapak/File/src/logger"
	"github.com/TimDebug/TutupLapak/File/src/models"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

var (
	redisOnce     sync.Once
	redisInstance *RedisClient
	appConfig     *config.Configuration = config.GetConfig()
)

func NewRedisClient() *RedisClient {
	redisOnce.Do(func() {
		redisInstance = &RedisClient{
			client: redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf(
					"%s:%s", appConfig.RedisHost, appConfig.RedisPort,
				),
				Password: "",
				DB:       0,
			}),
		}
	})

	// ping the client
	ctx := context.Background()
	_, err := redisInstance.client.Ping(ctx).Result()
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to connect to redis")
	}

	return redisInstance
}

func (c *RedisClient) Set(ctx context.Context, entity *models.FileEntity) error {
	// map file entity to map[string]string
	contents := map[string]string{
		"fileUri":          entity.FileURI,
		"fileThumbnailUri": entity.ThumbnailURI,
	}
	// serialize the contents of file entity
	data, err := json.Marshal(contents)
	if err != nil {
		return err
	}

	// store with key format: "file:{fileId}"
	key := "file:" + entity.FileID
	err = c.client.Set(ctx, key, string(data), 0).Err()
	if err != nil {
		return err
	}

	return nil
}
