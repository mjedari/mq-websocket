package hub

import (
	"errors"
	"github.com/google/uuid"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"repo.abanicon.com/abantheter-microservices/websocket/infra/utils"
)

type Hub struct {
	PublicRooms     *utils.SafeMap
	PrivateRooms    *utils.SafeMap
	ClientRooms     *utils.SafeMap
	PrivateReceiver chan PrivateMessage
	PublicReceiver  chan PublicMessage
	AuthReceiver    chan AuthMessage
}

func NewHub() *Hub {
	return &Hub{
		PrivateRooms:    utils.NewSafeMap(),
		PublicRooms:     utils.NewSafeMap(),
		ClientRooms:     utils.NewSafeMap(),
		PrivateReceiver: make(chan PrivateMessage),
		PublicReceiver:  make(chan PublicMessage),
		AuthReceiver:    make(chan AuthMessage),
	}
}

type RoomFactory func(name string) (contracts.IRoom, error) // todo: use this

func (h *Hub) GetPrivateRoom(name string, factory func(name string) (contracts.IPrivateRoom, error)) (contracts.IPrivateRoom, error) {
	if factory == nil {
		r, ok := h.PrivateRooms.Load(name)
		if !ok {
			return nil, errors.New("no room found")
		}
		return r.(contracts.IPrivateRoom), nil
	}

	newRoom, err := factory(name)
	if err != nil {
		return nil, err
	}

	r, _ := h.PrivateRooms.LoadOrStore(name, newRoom)
	return r.(contracts.IPrivateRoom), nil
}

func (h *Hub) GetPublicRoom(name string, factory func(name string) (contracts.IPublicRoom, error)) (contracts.IPublicRoom, error) {
	if factory == nil {
		r, ok := h.PublicRooms.Load(name)
		if !ok {
			return nil, errors.New("no room found")
		}
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
	h.ClientRooms.Store(clientId.String(), r)
}

func (h *Hub) GetClientRoom(clientId uuid.UUID) (contracts.IRoom, error) {
	r, ok := h.ClientRooms.Load(clientId.String())
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
	h.ClientRooms.Delete(clientId.String())
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
