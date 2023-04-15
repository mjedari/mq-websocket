package wsHandler

import "net/http"

type ChannelHandler interface {
	*http.Handler
	GetChannelType() string
	WriteMessage(message string)
}
