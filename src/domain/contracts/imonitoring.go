package contracts

type IMonitoring interface {
	AddClientToRoom(room string)
	RemoveClientFromRoom(room string)
	AuthenticationFailed()
	MessageReceived(room, kind string)
	AddConnection(kind string)
	RemoveConnection(kind string)
}
