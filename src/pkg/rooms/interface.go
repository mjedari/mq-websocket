package rooms

import "sync"

// maybe there whould be some other types of rooms?

type IRoom interface {
	GetName() string
	Broadcast(message []byte)
	PrivateSend(userId string, message []byte)
	GetClients() *sync.Map
	Leave(c *Client)
}

type IClient interface {
}
