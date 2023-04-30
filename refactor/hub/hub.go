package hub

import (
	"fmt"
	"sync"
	"websocket/refactor/room"
)

type Hub struct {
	rooms sync.Map
}

func NewHub() *Hub {
	return &Hub{}
}

func (h *Hub) GetRoom(name string, filter room.MessageFilter) room.IRoom {
	switch filter {
	case nil:
		newRoom := room.NewRoom(name)
		r, _ := h.rooms.LoadOrStore(name, newRoom)
		return r.(room.IRoom)
	default:
		newRoom := room.NewFilteredRoom(name, filter)
		r, _ := h.rooms.LoadOrStore(name, newRoom)
		return r.(room.IRoom)
	}
}

type PrivateMessage struct {
	UserId  []byte
	Room    []byte
	Message []byte
}

func (h *Hub) Streaming(privateChan chan PrivateMessage) {
	for {
		msg := <-privateChan
		fmt.Println("received private message")
		fmt.Println("room-id ", string(msg.Room))
		fmt.Println("user-id: ", string(msg.UserId))
		fmt.Println("message: ", string(msg.Message))

		h.rooms.Range(func(key, value any) bool {
			r, ok := value.(room.IRoom)
			if !ok {
				return false
			}
			if r.GetName() == string(msg.Room) {
				r.PrivateSend(string(msg.UserId), msg.Message)
			}

			return true
		})

	}
}
