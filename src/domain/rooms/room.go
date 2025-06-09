package rooms

import (
	"fmt"
	"github.com/mjedari/mq-websocket/domain/contracts"
	"github.com/mjedari/mq-websocket/infra/utils"
)

type BaseRoom struct {
	Name    string
	clients *utils.SafeMap
}

func NewBaseRoom(name string) (*BaseRoom, error) {
	// you can set some rules to prevent get new rooms by returning err
	return &BaseRoom{Name: name, clients: utils.NewSafeMap()}, nil
}

func (r *BaseRoom) GetName() string {
	return r.Name
}

func (r *BaseRoom) GetClients() *utils.SafeMap {
	return r.clients
}

func (r *BaseRoom) Leave(client contracts.IClient) bool {
	fmt.Printf("user \"%v\" is leaving \"%v\" rooms \n", client.GetId().String(), r.Name) //todo: fix this

	// todo: find out does commenting below line make memory leak or not
	//client.Leave()

	// remove client from room list
	_, existed := r.clients.LoadAndDelete(client.GetId().String())
	return existed
}

func (r *BaseRoom) SetClient(client contracts.IClient) {
	// this structure is storing inside map:
	// private: &{{e0acb902-1381-11ee-9651-363197453099 0x1400014a160 0x140000b25a0 0x140000b2600} 5 64802b08935acb1b0fa21e7f}
	// public: &{{f069daf0-1381-11ee-9651-363197453099 0x1400026e160 0x140002326c0 0x14000232720}}
	r.clients.Store(client.GetId().String(), client)
}

func (r *BaseRoom) Members() int {
	return r.clients.Len()
}
