package contracts

import (
	"sync"
)

type IRoom interface {
	GetName() string
	Broadcast(message []byte)
	PrivateSend(userId string, message []byte)
	GetClients() *sync.Map
	Leave(c IClient)
	SetClient(client IClient)
}
