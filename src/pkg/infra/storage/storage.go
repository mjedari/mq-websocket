package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"repo.abanicon.com/abantheter-microservices/websocket/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/utils"
	"time"
)

type Redis struct {
	*redis.Client
}

func NewRedis(conf configs.RedisConfig) (*redis.Client, error) {
	ctx := context.TODO()

	redisRetry, err := utils.Retry(func(ctx context.Context) (any, error) {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", conf.Host, conf.Port),
			Username: conf.User,
			Password: conf.Pass,
		})

		_, err := client.Ping(ctx).Result()
		if err != nil {
			return nil, err
		}

		return client, nil
	}, utils.RetryTimes, utils.RetryDelay)(ctx)

	if err != nil {
		return nil, err
	}
	// here we convert interface datatype to redis.Client
	client := redisRetry.(*redis.Client)

	return client, nil
}

func (r Redis) Store(ctx context.Context, key, value string, timeToLive time.Duration) error {
	err := r.Set(ctx, key, value, timeToLive).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r Redis) Fetch(ctx context.Context, key string) []byte {
	val, err := r.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	fmt.Println("fetch", val)
	return []byte(val)
}

func (r Redis) Exists(ctx context.Context, key string) bool {
	_, err := r.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	}
	return true
}

func (r Redis) Delete(ctx context.Context, key string) error {
	return r.Del(ctx, key).Err()
}
