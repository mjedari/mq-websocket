package hub

import "context"

type IKafkaHandler interface {
	Consume(ctx context.Context, topic string, responseChan chan KafkaMessage, privateChan chan PrivateMessage)
	CreateTopic(ctx context.Context, topics []string, partitions, replicationFactor int) error
}
