package hub

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"repo.abanicon.com/abantheter-microservices/websocket/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/rooms"
	"sync"
)

type Hub struct {
	kafka           IKafkaHandler
	rooms           sync.Map
	clientRooms     sync.Map
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

type RoomFactory func(name string) (rooms.IRoom, error)

func (h *Hub) GetRoom(name string, factory RoomFactory) (rooms.IRoom, error) {
	newRoom, err := factory(name)
	if err != nil {
		return nil, err
	}

	r, _ := h.rooms.LoadOrStore(name, newRoom)
	return r.(rooms.IRoom), nil
}

func (h *Hub) SetClientRoom(clientId uuid.UUID, r rooms.IRoom) error {
	// todo: investigating which one is better here: storing instance or its pointer?
	h.clientRooms.Store(clientId, r)
	return nil
}

func (h *Hub) GetClientRoom(clientId uuid.UUID) (rooms.IRoom, error) {
	r, ok := h.clientRooms.Load(clientId)
	if !ok {
		return nil, errors.New("not found")
	}

	if r == nil {
		return nil, nil
	}
	rrr := r.(rooms.IRoom)
	return rrr, nil
}

func (h *Hub) RemoveClientRoom(clientId uuid.UUID) error {
	h.clientRooms.Delete(clientId)
	return nil
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
		fmt.Println("rooms-id ", string(msg.Room))
		fmt.Println("user-id: ", string(msg.UserId))
		fmt.Println("message: ", string(msg.Message))

		h.rooms.Range(func(key, value any) bool {
			r, ok := value.(rooms.IRoom)
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
	go h.kafka.Consume(context.Background(), configs.Config.Topics.PublicTopic, h.PublicReceiver, h.PrivateReceiver)

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
		fmt.Println("rooms-id ", "public")
		fmt.Println("message: ", msg.Value)

		r, ok := h.rooms.Load("public")
		if !ok {
			fmt.Println("not found public channel")
			continue
		}

		publicRoom := r.(rooms.IRoom)
		publicRoom.Broadcast([]byte(msg.Value))
	}
}

func (h *Hub) LeaveClient(ctx context.Context, userId string) {
	h.rooms.Range(func(room, anyRoom any) bool {
		r, ok := anyRoom.(rooms.IRoom)
		if !ok {
			return false
		}

		r.GetClients().Range(func(client, status any) bool {
			c := client.(*rooms.Client)

			if c.UserId == userId {
				r.Leave(c)
			}
			return true
		})

		return true
	})
}
