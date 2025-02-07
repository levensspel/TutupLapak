package cache

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/response"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/service"
	"github.com/bytedance/sonic"
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

	// Set many user profiles
	MSetUserProfiles(ctx context.Context, users []response.UserWithIdResponse) error

	GetUserProfile(ctx context.Context, userId string) (*response.UserResponse, bool)
	MGetUserProfiles(ctx context.Context, keys []string) (cachedUsers []response.UserWithIdResponse, missedUserIds []string, ok bool)
	GetFile(ctx context.Context, fileId string) (*service.File, bool)
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

	return d.setAsMap(ctx, "user:"+userId, userMap)
}

// Utilizing the Redis `MSET` command.
// Sets the given keys to their respective values.
// MSET replaces existing values with new values,
// just as regular SET. See MSETNX if you don't
// want to overwrite existing values. MSET is atomic,
// so all given keys are set at once. It is not possible
// for clients to see that some of the keys were updated while others are unchanged.
// Src: https://redis.io/docs/latest/commands/mset/
func (d RedisCacheClient) MSetUserProfiles(ctx context.Context, users []response.UserWithIdResponse) error {
	usersMap := make(map[string]map[string]string)
	for _, user := range users {
		usersMap[user.UserId] = map[string]string{
			"email":             user.Email,
			"phone":             user.Phone,
			"fileId":            user.FileId,
			"fileUri":           user.FileUri,
			"fileThumbnailUri":  user.FileThumbnailUri,
			"bankAccountName":   user.BankAccountName,
			"bankAccountHolder": user.BankAccountHolder,
			"bankAccountNumber": user.BankAccountNumber,
		}
	}

	cacheMap := make(map[string]string)
	for k, v := range usersMap {
		data, err := sonic.Marshal(v)
		if err != nil {
			return err
		}

		cacheMap["user:"+k] = string(data)
	}

	_, err := d.client.MSet(ctx, cacheMap).Result()
	if err != nil {
		return err
	}
	return nil
}

func (d RedisCacheClient) GetUserProfile(ctx context.Context, userId string) (*response.UserResponse, bool) {
	defer func() {
		recover()
	}()

	result, ok := d.getAsMap(ctx, "user:"+userId)
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

// Get many user profiles by utilizing the Redis' `MGET` command to query once to Redis server
func (d RedisCacheClient) MGetUserProfiles(ctx context.Context, keys []string) (cachedUsers []response.UserWithIdResponse, missedUserIds []string, ok bool) {
	val, err := d.client.MGet(ctx, keys...).Result()
	if err != nil {
		return []response.UserWithIdResponse{}, []string{}, false
	}

	cachedMap := make(map[string]map[string]string)

	count := 0
	for k, v := range val {
		if vStr, ok := v.(string); ok {
			var r map[string]string
			err = sonic.Unmarshal([]byte(vStr), &r)
			if err == nil {
				count++
				cachedMap[keys[k]] = r
			}
		} else {
			cachedMap[keys[k]] = map[string]string{}
		}
	}
	if count == 0 {
		return []response.UserWithIdResponse{}, []string{}, false
	}

	for key, userCache := range cachedMap {
		id := key[5:] // because the cache key format "user:SOME_ID"

		if len(userCache) == 0 {
			missedUserIds = append(missedUserIds, id)
		} else {
			user := response.UserWithIdResponse{
				UserId:            id,
				Email:             userCache["email"],
				Phone:             userCache["phone"],
				FileId:            userCache["fileId"],
				FileUri:           userCache["fileUri"],
				FileThumbnailUri:  userCache["fileThumbnailUri"],
				BankAccountName:   userCache["bankAccountName"],
				BankAccountHolder: userCache["bankAccountHolder"],
				BankAccountNumber: userCache["bankAccountNumber"],
			}

			cachedUsers = append(cachedUsers, user)
		}
	}
	return cachedUsers, missedUserIds, true
}

func (d RedisCacheClient) GetFile(ctx context.Context, fileId string) (*service.File, bool) {
	defer func() {
		recover()
	}()

	result, ok := d.getAsMap(ctx, "file:"+fileId)
	if !ok {
		return &service.File{}, false
	}

	return &service.File{
		FileUri:          result["fileUri"],
		FileThumbnailUri: result["fileThumbnailUri"],
	}, true
}
