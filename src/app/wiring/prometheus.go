package wiring

import "repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"

func (w *Wire) GetMonitoringService() contracts.IMonitoring {
	return w.Monitoring
}
