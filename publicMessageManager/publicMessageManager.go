package publicMessageManager

import (
	"context"
	"websocket/configs"
	"websocket/kafkaManager"
	"websocket/wsHandler"
)

func ReceiveMessages() {
	messagesChan := make(chan kafkaManager.KafkaMessage)
	go kafkaManager.Consume(context.Background(), configs.WebSocketPublicTopic, messagesChan)

	for {
		resp := <-messagesChan
		wsHandler.PublicWriter(resp.Value)
	}
}
