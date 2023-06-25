package contracts

import "github.com/google/uuid"

//type RoomFactory func(name string) (IRoom, error)

type IHub interface {
	GetRoom(name string, factory func(name string) (IRoom, error)) (IRoom, error)
	SetClientRoom(clientId uuid.UUID, r IRoom)
	GetClientRoom(clientId uuid.UUID) (IRoom, error)
	RemoveClientRoom(clientId uuid.UUID)
}
