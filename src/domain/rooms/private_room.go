package rooms

import (
	"fmt"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
)

type PrivateRoom struct {
	*BaseRoom
}

func NewPrivateRoom(name string) (*PrivateRoom, error) {
	// you can set some rules to prevent get new rooms by returning err
	return &PrivateRoom{&BaseRoom{
		Name: name,
	}}, nil
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
