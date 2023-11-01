package contracts

import (
	"sync"
)

type IRoom interface {
	GetName() string
	GetClients() *sync.Map
	Leave(c IClient)
	SetClient(client IClient)
}

type IPublicRoom interface {
	IRoom
	Broadcast(message []byte)
}

type IPrivateRoom interface {
	IRoom
	PrivateSend(userId string, message []byte)
}
