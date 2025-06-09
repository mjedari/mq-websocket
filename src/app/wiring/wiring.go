package wiring

import (
	"github.com/mjedari/mq-websocket/app/configs"
	"github.com/mjedari/mq-websocket/domain/contracts"
	"github.com/mjedari/mq-websocket/domain/hub"
	"github.com/mjedari/mq-websocket/infra/broker"
	"github.com/mjedari/mq-websocket/infra/monitoring"
	"github.com/mjedari/mq-websocket/infra/rate_limiter"
)

var Wiring *Wire

type Wire struct {
	Kafka       *broker.Kafka
	Redis       contracts.IStorage
	Monitoring  contracts.IMonitoring
	RateLimiter *rate_limiter.RateLimiter
	Hub         *hub.Hub
	Configs     configs.Configuration
}

func NewWire(kafka *broker.Kafka, redis contracts.IStorage, limiter *rate_limiter.RateLimiter, configs configs.Configuration) *Wire {
	newHub := hub.NewHub()
	newMonitoring := monitoring.NewPrometheus()
	return &Wire{Kafka: kafka, Redis: redis, Monitoring: newMonitoring, Hub: newHub, RateLimiter: limiter, Configs: configs}
}

func (w *Wire) GetRateLimiter() *rate_limiter.RateLimiter {
	return w.RateLimiter
}

func (w *Wire) GetStorage() contracts.IStorage {
	return w.Redis
}
