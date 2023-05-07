package hub

import (
	"context"
	"fmt"
	"sync"
	"websocket/configs"
	"websocket/refactor/room"
)

type Hub struct {
	kafka           IKafkaHandler
	rooms           sync.Map
	PrivateReceiver chan PrivateMessage
	PublicReceiver  chan KafkaMessage
	AuthReceiver    chan KafkaMessage
}

func NewHub(kafka IKafkaHandler) *Hub {
	return &Hub{
		kafka:           kafka,
		PrivateReceiver: make(chan PrivateMessage),
		PublicReceiver:  make(chan KafkaMessage),
		AuthReceiver:    make(chan KafkaMessage),
	}
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

type KafkaMessage struct {
	Value         string
	CorrelationId string
}

func (h *Hub) PrivateStreaming() {
	for {
		msg := <-h.PrivateReceiver
		//msg := <-privateChan
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

func (h *Hub) Streaming() {
	go h.kafka.Consume(context.Background(), configs.WebSocketPublicTopic, h.PublicReceiver, h.PrivateReceiver)
	go h.kafka.Consume(context.Background(), configs.WEBSOCKET_AUTHENTICATION_TOPIC, h.AuthReceiver, nil)

	// listening for private channel
	go h.PrivateStreaming()

	// listening for public channel
	go h.PublicStreaming()

}

func (h *Hub) PublicStreaming() {
	for {
		msg := <-h.PublicReceiver
		//wsHandler.PublicWriter(resp.Value)

		fmt.Println("received public message")
		fmt.Println("room-id ", "public")
		fmt.Println("message: ", msg.Value)

		r, ok := h.rooms.Load("public")
		if !ok {
			fmt.Println("not found public channel")
			continue
		}

		publicRoom := r.(room.IRoom)
		publicRoom.Broadcast([]byte(msg.Value))
	}
}
