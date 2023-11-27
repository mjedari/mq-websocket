package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Prometheus struct {
	clientsInRoom        *prometheus.GaugeVec
	messagesReceived     *prometheus.CounterVec
	authenticationFailed prometheus.Counter
}

func NewPrometheus() *Prometheus {
	clientsInRoom := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "clients_in_room",
			Help: "Number of clients in each room",
		},
		[]string{"room_name"},
	)

	messagesReceived := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_received",
			Help: "Number of messages receiving",
		},
		[]string{"room", "type"},
	)

	authenticationFailed := promauto.NewCounter(prometheus.CounterOpts{
		Name: "failed_authentication_total",
		Help: "Total number of failed transactions",
	})

	return &Prometheus{clientsInRoom: clientsInRoom, messagesReceived: messagesReceived, authenticationFailed: authenticationFailed}
}

func (p Prometheus) AddClientToRoom(room string) {
	p.clientsInRoom.With(prometheus.Labels{"room_name": room}).Inc()
}

func (p Prometheus) RemoveClientFromRoom(room string) {
	p.clientsInRoom.With(prometheus.Labels{"room_name": room}).Dec()
}

func (p Prometheus) AuthenticationFailed() {
	p.authenticationFailed.Inc()
}

func (p Prometheus) MessageReceived(room, kind string) {
	p.messagesReceived.With(prometheus.Labels{"room": room, "type": kind}).Inc()
}
