package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"repo.abanicon.com/abantheter-microservices/websocket/configs"
)

type Redis struct {
	*redis.Client
}

func NewRedis(conf configs.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		Username: conf.User,
		Password: conf.Pass,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (redis Redis) StoreUntil() {

}
