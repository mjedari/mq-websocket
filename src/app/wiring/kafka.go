package wiring

import "repo.abanicon.com/abantheter-microservices/websocket/infra/broker"

func (w *Wire) GetKafkaAdmin() *broker.Kafka {
	return w.Kafka
}
