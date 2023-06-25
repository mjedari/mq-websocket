package hub

import (
	"errors"
	"github.com/google/uuid"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"sync"
)

type Hub struct {
	Rooms           sync.Map
	clientRooms     sync.Map
	PrivateReceiver chan PrivateMessage
	PublicReceiver  chan PublicMessage
	AuthReceiver    chan AuthMessage
}

func NewHub() *Hub {
	return &Hub{
		PrivateReceiver: make(chan PrivateMessage),
		PublicReceiver:  make(chan PublicMessage),
		AuthReceiver:    make(chan AuthMessage),
	}
}

type RoomFactory func(name string) (contracts.IRoom, error) // todo: use this

func (h *Hub) GetRoom(name string, factory func(name string) (contracts.IRoom, error)) (contracts.IRoom, error) {
	newRoom, err := factory(name)
	if err != nil {
		return nil, err
	}

	r, _ := h.Rooms.LoadOrStore(name, newRoom)
	return r.(contracts.IRoom), nil
}

func (h *Hub) SetClientRoom(clientId uuid.UUID, r contracts.IRoom) {
	// todo: investigating which one is better here: storing instance or its pointer?
	h.clientRooms.Store(clientId, r)
}

func (h *Hub) GetClientRoom(clientId uuid.UUID) (contracts.IRoom, error) {
	r, ok := h.clientRooms.Load(clientId)
	if !ok {
		return nil, errors.New("not found")
	}

	if r == nil {
		return nil, nil
	}
	rrr := r.(contracts.IRoom)
	return rrr, nil
}

func (h *Hub) RemoveClientRoom(clientId uuid.UUID) {
	h.clientRooms.Delete(clientId)
}

// todo: remove these

type PrivateMessage struct {
	UserId  []byte
	Room    []byte
	Message []byte
}

type PublicMessage struct {
	Value         string
	CorrelationId string
}

type AuthMessage struct {
	UserId      string
	DeviceId    string
	AccessToken string
}
