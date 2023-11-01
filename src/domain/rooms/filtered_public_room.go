package rooms

type MessageFilter func(message []byte) bool

type FilteredPublicRoom struct {
	*PublicRoom
	filter MessageFilter
}

func NewFilteredPublicRoom(name string, filter MessageFilter) *FilteredPublicRoom {
	// you can set some rules to prevent get new rooms by returning err

	return &FilteredPublicRoom{
		PublicRoom: &PublicRoom{
			BaseRoom: &BaseRoom{
				Name: name,
			},
		},
		filter: filter,
	}
}

func (fr *FilteredPublicRoom) Broadcast(message []byte) {
	if fr.filter == nil || fr.filter(message) {
		fr.PublicRoom.Broadcast(message)
	}
}
