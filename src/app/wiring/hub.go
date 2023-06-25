package wiring

import "repo.abanicon.com/abantheter-microservices/websocket/domain/hub"

func (w *Wire) GetHub() *hub.Hub {
	return w.Hub
}
