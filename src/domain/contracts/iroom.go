package contracts

import (
	"repo.abanicon.com/abantheter-microservices/websocket/infra/utils"
)

type IRoom interface {
	GetName() string
	GetClients() *utils.SafeMap
	Leave(c IClient) bool
	SetClient(client IClient)
	Members() int
}

type IPublicRoom interface {
	IRoom
	Broadcast(message []byte)
}

type IPrivateRoom interface {
	IRoom
	PrivateSend(userId string, message []byte)
}
