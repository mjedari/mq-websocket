package contracts

import (
	"context"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/broker"
)

type IBroker interface {
	ConsumeAuth(ctx context.Context, topic string, responseChan chan broker.ResponseMessage)
}
