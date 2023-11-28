package contracts

import "github.com/google/uuid"

//type RoomFactory func(name string) (IRoom, error)

type IHub interface {
	GetPublicRoom(name string, factory func(name string) (IPublicRoom, error)) (IPublicRoom, error)
	GetPrivateRoom(name string, factory func(name string) (IPrivateRoom, error)) (IPrivateRoom, error)
	SetClientRoom(clientId uuid.UUID, r IRoom)
	GetClientRoom(clientId uuid.UUID) (IRoom, error)
	RemoveClientRoom(clientId uuid.UUID)
}
