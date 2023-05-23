package rooms

import (
	"fmt"
	"sync"
)

type Room struct {
	Name    string
	Clients sync.Map
	mux     sync.Mutex
}

func NewRoom(name string) (*Room, error) {
	// you can set some rules to prevent get new room by returning err
	return &Room{Name: name}, nil
}

func (r *Room) GetName() string {
	return r.Name
}

func (r *Room) GetClients() *sync.Map {
	return &r.Clients
}

func (r *Room) Broadcast(message []byte) {
	//r.mux.Lock()
	//defer r.mux.Unlock()
	var clientNumbers uint64

	r.Clients.Range(func(key, value any) bool {
		clientNumbers++

		c, ok := key.(*Client)
		if !ok {
			return false
		}
		fmt.Println("published by streaming on client:", c.UserId)
		c.Send <- message
		return true
	})
	fmt.Printf("published to #%v clients\n", clientNumbers)

}

func (r *Room) PrivateSend(userId string, message []byte) {
	//r.mux.Lock()
	//defer r.mux.Unlock()

	r.Clients.Range(func(key, value any) bool {
		c, ok := key.(*Client)
		if !ok {
			return false
		}

		if c.UserId == userId {
			fmt.Println("published by streaming on client privately:", c.UserId)
			c.Send <- message
		}
		return true
	})

}

func (r *Room) Leave(c *Client) {
	fmt.Printf("user \"%v\" is leaving \"%v\" rooms \n", c.UserId, r.Name)
	//r.mux.Lock()
	//defer r.mux.Unlock()

	c.Close <- true
	r.Clients.Delete(c)

}
