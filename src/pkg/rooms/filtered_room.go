package rooms

import "sync"

type MessageFilter func(message []byte) bool

type FilteredRoom struct {
	*Room
	filter MessageFilter
}

func (fr *FilteredRoom) GetClients() *sync.Map {
	return &fr.Clients
}

func NewFilteredRoom(name string, filter MessageFilter) *FilteredRoom {
	// you can set some rules to prevent get new room by returning err

	return &FilteredRoom{
		Room: &Room{
			Name: name,
		},
		filter: filter,
	}
}

func (fr *FilteredRoom) Broadcast(message []byte) {
	if fr.filter == nil || fr.filter(message) {
		fr.Room.Broadcast(message)
	}
}
