package contracts

type IMonitoring interface {
	AddClientToRoom(room string)
	RemoveClientFromRoom(room string)
	AuthenticationFailed()
}
