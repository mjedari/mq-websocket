package publicMessageManager

import (
	"context"
	"websocket/configs"
	"websocket/kafkaManager"
	"websocket/refactor/hub"
	"websocket/wsHandler"
)

func ReceiveMessages(privateChan chan hub.PrivateMessage) {
	messagesChan := make(chan kafkaManager.KafkaMessage)
	go kafkaManager.Consume(context.Background(), configs.WebSocketPublicTopic, messagesChan, privateChan)

	for {
		resp := <-messagesChan
		wsHandler.PublicWriter(resp.Value)
	}
}
