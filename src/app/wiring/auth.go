package wiring

import (
	"repo.abanicon.com/abantheter-microservices/websocket/app/messaging"
)

func (w *Wire) GetAuthService() *messaging.AuthService {
	return messaging.NewAuthService(w.GetStorage(), w.GetMonitoringService(), w.GetHub())
}
