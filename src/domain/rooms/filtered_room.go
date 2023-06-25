package rooms

type MessageFilter func(message []byte) bool

type FilteredRoom struct {
	*BaseRoom
	filter MessageFilter
}

func NewFilteredRoom(name string, filter MessageFilter) *FilteredRoom {
	// you can set some rules to prevent get new rooms by returning err

	return &FilteredRoom{
		BaseRoom: &BaseRoom{
			Name: name,
		},
		filter: filter,
	}
}

func (fr *FilteredRoom) Broadcast(message []byte) {
	if fr.filter == nil || fr.filter(message) {
		fr.BaseRoom.Broadcast(message)
	}
}
