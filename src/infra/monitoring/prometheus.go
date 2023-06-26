package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Prometheus struct {
	clientsInRoom *prometheus.GaugeVec
}

func NewPrometheus() *Prometheus {
	clientsInRoom := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "clients_in_room",
			Help: "Number of clients in each room",
		},
		[]string{"room_name"},
	)

	return &Prometheus{clientsInRoom: clientsInRoom}
}

func (p Prometheus) AddClientToRoom(room string) {
	p.clientsInRoom.With(prometheus.Labels{"room_name": room}).Inc()
}

func (p Prometheus) RemoveClientFromRoom(room string) {
	p.clientsInRoom.With(prometheus.Labels{"room_name": room}).Dec()
}
