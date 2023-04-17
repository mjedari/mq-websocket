package wsHandler

import "net/http"

type ChannelHandler interface {
	*http.Handler
	GetBrokerTopic() string
	WriteMessage(message string)
}
