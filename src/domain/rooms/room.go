package rooms

import (
	"fmt"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"sync"
)

type BaseRoom struct {
	Name    string
	Clients sync.Map
}

func NewRoom(name string) (*BaseRoom, error) {
	// you can set some rules to prevent get new rooms by returning err
	return &BaseRoom{Name: name}, nil
}

func (r *BaseRoom) GetName() string {
	return r.Name
}

func (r *BaseRoom) GetClients() *sync.Map {
	return &r.Clients
}

func (r *BaseRoom) Broadcast(message []byte) {
	var clientNumbers uint64

	r.Clients.Range(func(key, value any) bool {
		clientNumbers++

		c, ok := value.(contracts.IClient)
		if !ok {
			return false
		}
		//fmt.Println("published by streaming on clients:", c.Id)
		c.SendMessage(message)
		return true
	})
	fmt.Printf("broadcasted to #%v clients\n", clientNumbers)
}

func (r *BaseRoom) PrivateSend(userId string, message []byte) {
	var clientNumbers uint64
	r.Clients.Range(func(key, value any) bool {
		clientNumbers++
		c, ok := value.(contracts.IClient)
		if !ok {
			return false
		}

		if c.Check(userId) {
			//fmt.Println("published by streaming on clients privately:", c.UserId)
			c.SendMessage(message)
		}

		return true
	})
	fmt.Printf("published to #%v private clients\n", clientNumbers)
}

func (r *BaseRoom) Leave(client contracts.IClient) {
	fmt.Printf("user \"%v\" is leaving \"%v\" rooms \n", "", r.Name) //todo: fix this

	client.Leave()

	// remove client from room list
	r.Clients.Delete(client.GetId())

}

func (r *BaseRoom) SetClient(client contracts.IClient) {
	// this structure is storing inside map:
	// private: &{{e0acb902-1381-11ee-9651-363197453099 0x1400014a160 0x140000b25a0 0x140000b2600} 5 64802b08935acb1b0fa21e7f}
	// public: &{{f069daf0-1381-11ee-9651-363197453099 0x1400026e160 0x140002326c0 0x14000232720}}
	r.Clients.Store(client.GetId(), client)
}
