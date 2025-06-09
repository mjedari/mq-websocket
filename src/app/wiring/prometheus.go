package wiring

import "github.com/mjedari/mq-websocket/domain/contracts"

func (w *Wire) GetMonitoringService() contracts.IMonitoring {
	return w.Monitoring
}
