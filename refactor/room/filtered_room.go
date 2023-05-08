package room

type MessageFilter func(message []byte) bool

type FilteredRoom struct {
	*Room
	filter MessageFilter
}

func (fr *FilteredRoom) GetClients() map[*Client]bool {
	return fr.Clients
}

func NewFilteredRoom(name string, filter MessageFilter) *FilteredRoom {
	return &FilteredRoom{
		Room: &Room{
			Name:    name,
			Clients: make(map[*Client]bool),
		},
		filter: filter,
	}
}

func (fr *FilteredRoom) Broadcast(message []byte) {
	if fr.filter == nil || fr.filter(message) {
		fr.Room.Broadcast(message)
	}
}
