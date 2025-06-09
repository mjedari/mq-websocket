package wiring

import "github.com/mjedari/mq-websocket/domain/hub"

func (w *Wire) GetHub() *hub.Hub {
	return w.Hub
}
