package wiring

import (
	"repo.abanicon.com/abantheter-microservices/websocket/app/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/infra/broker"
	"repo.abanicon.com/abantheter-microservices/websocket/infra/rate_limiter"
)

var Wiring *Wire

type Wire struct {
	Kafka       *broker.Kafka
	Redis       contracts.IStorage
	RateLimiter *rate_limiter.RateLimiter
	Hub         *hub.Hub
	Configs     configs.Configuration
}

func NewWire(kafka *broker.Kafka, redis contracts.IStorage, limiter *rate_limiter.RateLimiter, configs configs.Configuration) *Wire {
	newHub := hub.NewHub()
	return &Wire{Kafka: kafka, Redis: redis, Hub: newHub, RateLimiter: limiter, Configs: configs}
}

func (w *Wire) GetRateLimiter() *rate_limiter.RateLimiter {
	return w.RateLimiter
}

func (w *Wire) GetStorage() contracts.IStorage {
	return w.Redis
}
