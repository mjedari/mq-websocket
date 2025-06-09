package wiring

import (
	"github.com/mjedari/mq-websocket/app/messaging"
)

func (w *Wire) GetAuthService() *messaging.AuthService {
	return messaging.NewAuthService(w.GetStorage(), w.GetMonitoringService(), w.GetHub())
}
