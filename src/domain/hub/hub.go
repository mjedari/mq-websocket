package hub

import (
	"errors"
	"github.com/google/uuid"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"sync"
)

type Hub struct {
	PublicRooms     sync.Map
	PrivateRooms    sync.Map
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

func (h *Hub) GetPrivateRoom(name string, factory func(name string) (contracts.IPrivateRoom, error)) (contracts.IPrivateRoom, error) {
	newRoom, err := factory(name)
	if err != nil {
		return nil, err
	}

	r, _ := h.PrivateRooms.LoadOrStore(name, newRoom)
	return r.(contracts.IPrivateRoom), nil
}

func (h *Hub) GetPublicRoom(name string, factory func(name string) (contracts.IPublicRoom, error)) (contracts.IPublicRoom, error) {
	if factory == nil {
		r, _ := h.PublicRooms.Load(name)
		return r.(contracts.IPublicRoom), nil
	}

	newRoom, err := factory(name)
	if err != nil {
		return nil, err
	}

	r, _ := h.PublicRooms.LoadOrStore(name, newRoom)
	return r.(contracts.IPublicRoom), nil
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
	Room          []byte
	Message       []byte
	CorrelationId string
}

type AuthMessage struct {
	UserId      string
	DeviceId    string
	AccessToken string
}
