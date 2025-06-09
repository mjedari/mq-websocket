package rooms

import (
	"fmt"
	"github.com/mjedari/mq-websocket/domain/contracts"
)

type PublicRoom struct {
	*BaseRoom
}

func NewPublicRoom(name string) (*PublicRoom, error) {
	// you can set some rules to prevent get new rooms by returning err
	baseRoom, _ := NewBaseRoom(name)
	return &PublicRoom{baseRoom}, nil
}

func (r *PublicRoom) Broadcast(message []byte) {
	var clientNumbers uint64

	r.clients.Range(func(key, value any) bool {
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
