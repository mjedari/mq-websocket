package wiring

import (
	"github.com/go-redis/redis/v8"
	"repo.abanicon.com/abantheter-microservices/websocket/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/auth"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/storage"
)

var Wiring *Wire

type Wire struct {
	Kafka   hub.IKafkaHandler
	Redis   *redis.Client
	Configs configs.Configuration
}

func NewWire(kafka hub.IKafkaHandler, redis *redis.Client, configs configs.Configuration) *Wire {
	return &Wire{Kafka: kafka, Redis: redis, Configs: configs}
}

func (w *Wire) GetRedis() *storage.Redis {
	return &storage.Redis{Client: w.Redis}
}

func (w *Wire) GetKafkaAdmin() hub.IKafkaHandler {
	return w.Kafka
}

func (w *Wire) GetAuthService() *auth.AuthService {
	return auth.NewAuthService(w.GetRedis(), w.GetAuthServiceConfig())
}
