package cache

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/response"
	"github.com/samber/do/v2"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultCacheTtl = 5 * time.Minute
)

var (
	redisCacheClientOnce     sync.Once
	redisCacheClientInstance *RedisCacheClient
)

type CacheClientInterface interface {
	// Internal Interfaces
	set(ctx context.Context, key string, value string) error
	setAsMap(ctx context.Context, key string, value map[string]string) error
	get(ctx context.Context, key string) (string, bool)
	getAsMap(ctx context.Context, key string) (map[string]string, bool)
	getAsMapArray(ctx context.Context, key string) ([]map[string]string, bool)

	// Exposed Interfaces to Services
	SetUserProfile(ctx context.Context, userId string, user *response.UserResponse) error
	GetUserProfile(ctx context.Context, userId string) (*response.UserResponse, bool)
}

type RedisCacheClient struct {
	client *redis.Client
}

func NewRedisCacheClient() CacheClientInterface {
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

func NewRedisCacheClientInject(i do.Injector) (RedisCacheClient, error) {
	redisCacheClientInstance = &RedisCacheClient{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf(
				"%s:%s",
				config.GetRedisHost(),
				config.GetRedisPort(),
			),
			Password: config.GetRedisPassword(), // default no passoword
			DB:       config.GetRedisDbCount(),  // Default database number
		}),
	}

	return *redisCacheClientInstance, nil
}

func (d RedisCacheClient) SetUserProfile(ctx context.Context, userId string, user *response.UserResponse) error {
	defer func() {
		recover()
	}()

	userMap := map[string]string{
		"email":             user.Email,
		"phone":             user.Phone,
		"fileId":            user.FileId,
		"fileUri":           user.FileUri,
		"fileThumbnailUri":  user.FileThumbnailUri,
		"bankAccountName":   user.BankAccountName,
		"bankAccountHolder": user.BankAccountHolder,
		"bankAccountNumber": user.BankAccountNumber,
	}

	return d.setAsMap(ctx, userId, userMap)
}

func (d RedisCacheClient) GetUserProfile(ctx context.Context, userId string) (*response.UserResponse, bool) {
	defer func() {
		recover()
	}()

	result, ok := d.getAsMap(ctx, userId)
	if !ok {
		return &response.UserResponse{}, false
	}

	return &response.UserResponse{
		Email:             result["email"],
		Phone:             result["phone"],
		FileId:            result["fileId"],
		FileUri:           result["fileUri"],
		FileThumbnailUri:  result["fileThumbnailUri"],
		BankAccountName:   result["bankAccountName"],
		BankAccountHolder: result["bankAccountHolder"],
		BankAccountNumber: result["bankAccountNumber"],
	}, true
}
