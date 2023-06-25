package contracts

import (
	"context"
)

type IBroker interface {
	// todo: remove userid use general one like headers
	Consume(ctx context.Context, topic string, publicResponseFunction, privateResponseFunction func(userId, key, value []byte)) // how can I set this
	Produce(ctx context.Context, message IBrokerMessage)
	ConsumeAuth(ctx context.Context, topic string, authResponseFunction func(userId, key, value []byte)) // todo: remove
	CreateTopics(ctx context.Context, topics []string, partitions, replicationFactor int) error
	ConsumeHealth(ctx context.Context, topic string) ([]byte, error)
}

type IBrokerMessage interface {
	GetTopic() string
	GetKey() []byte
	GetMessage() []byte
	GetResponseTopic() string
	GetRequestId() string
}
