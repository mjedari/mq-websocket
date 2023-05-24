package wiring

import (
	"github.com/go-redis/redis/v8"
	"repo.abanicon.com/abantheter-microservices/websocket/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/auth"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/broker"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/infra/storage"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/rate_limiter"
)

var Wiring *Wire

type Wire struct {
	Kafka       *broker.Kafka
	Redis       *redis.Client
	RateLimiter *rate_limiter.RateLimiter
	Hub         *hub.Hub
	Configs     configs.Configuration
}

func NewWire(kafka *broker.Kafka, redis *redis.Client, hub *hub.Hub, limiter *rate_limiter.RateLimiter, configs configs.Configuration) *Wire {
	return &Wire{Kafka: kafka, Redis: redis, Hub: hub, RateLimiter: limiter, Configs: configs}
}

func (w *Wire) GetRedis() *storage.Redis {
	return &storage.Redis{Client: w.Redis}
}

func (w *Wire) SetNewRedisInstance() error {
	newInstance, err := storage.NewRedis(w.Configs.Redis)
	if err != nil {
		return err
	}
	w.Redis = newInstance
	return nil
}

func (w *Wire) GetKafkaAdmin() *broker.Kafka {
	return w.Kafka
}

func (w *Wire) GetNewAuthService() *auth.Auth {
	return auth.NewAuth(w.GetRedis(), w.Kafka, w.GetHub())
}

func (w *Wire) GetAuthService() *auth.AuthService {
	return auth.NewAuthService(w.GetRedis(), w.GetAuthServiceConfig())
}

func (w *Wire) GetHub() *hub.Hub {
	return w.Hub
}

func (w *Wire) GetRateLimiter() *rate_limiter.RateLimiter {
	return w.RateLimiter
}
