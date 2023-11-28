package rooms

import (
	"fmt"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
)

type PublicRoom struct {
	*BaseRoom
}

func NewPublicRoom(name string) (*PublicRoom, error) {
	// you can set some rules to prevent get new rooms by returning err
	return &PublicRoom{&BaseRoom{
		Name: name,
	}}, nil
}

func (r *PublicRoom) Broadcast(message []byte) {
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
