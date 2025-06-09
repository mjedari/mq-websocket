package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mjedari/mq-websocket/app/configs"
	"github.com/mjedari/mq-websocket/infra/utils"
	"time"
)

// todo: retry pattern for storage

type Redis struct {
	Client *redis.Client
	Config configs.RedisConfig
}

func (r Redis) CheckHealth(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r Redis) ResetConnection(ctx context.Context) error {
	newClient, err := NewRedis(r.Config)
	if err != nil {
		return err
	}
	r.Client = newClient.Client
	return nil
}

func NewRedis(conf configs.RedisConfig) (*Redis, error) {
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

	return &Redis{Client: client}, nil
}

func (r Redis) Store(ctx context.Context, key, value string, timeToLive time.Duration) error {
	err := r.Client.Set(ctx, key, value, timeToLive).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r Redis) Fetch(ctx context.Context, key string) []byte {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	return []byte(val)
}

func (r Redis) Exists(ctx context.Context, key string) bool {
	_, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	}
	return true
}

func (r Redis) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
