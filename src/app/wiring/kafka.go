package wiring

import "github.com/mjedari/mq-websocket/infra/broker"

func (w *Wire) GetKafkaAdmin() *broker.Kafka {
	return w.Kafka
}
