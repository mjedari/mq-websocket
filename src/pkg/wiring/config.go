package wiring

import (
	"net"
	"repo.abanicon.com/abantheter-microservices/websocket/configs"
)

func (w *Wire) GetServerConfig() string {
	return net.JoinHostPort(w.Configs.Server.Host, w.Configs.Server.Port)
}

func (w *Wire) GetAuthServiceConfig() configs.AuthServer {
	return w.Configs.AuthServer
}

func (w *Wire) GetAuthServerConfig() configs.AuthServer {
	return w.Configs.AuthServer
}

func (w *Wire) GetKafkaConfig() configs.KafkaConfig {
	return w.Configs.Kafka
}

func (w *Wire) GetRedisConfig() configs.RedisConfig {
	return w.Configs.Redis
}
