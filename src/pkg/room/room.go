package room

import (
	"fmt"
	"sync"
)

type Room struct {
	Name    string
	Clients map[*Client]bool
	mux     sync.Mutex
}

func NewRoom(name string) *Room {
	return &Room{Name: name, Clients: make(map[*Client]bool)}
}

func (r *Room) GetName() string {
	return r.Name
}

func (r *Room) GetClients() map[*Client]bool {
	return r.Clients
}

func (r *Room) Broadcast(message []byte) {
	//r.mux.Lock()
	//defer r.mux.Unlock()

	for client := range r.Clients {
		fmt.Println("published by streaming on client:", client)
		client.Send <- message

	}
}

func (r *Room) PrivateSend(userId string, message []byte) {
	//r.mux.Lock()
	//defer r.mux.Unlock()

	for client := range r.Clients {
		if client.UserId == userId {
			fmt.Println("published by streaming on client privately:", client)
			client.Send <- message
		}
	}
}

func (r *Room) Leave(c *Client) {
	fmt.Printf("user \"%v\" is leaving \"%v\" room \n", c.UserId, r.Name)
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.Clients, c)
}
