package wsHandler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
	"websocket/configs"
	"websocket/kafkaManager"
)

var privateConnections sync.Map

type PrivateHandler struct {
	//
}

func NewPrivateHandler() *PrivateHandler {
	return &PrivateHandler{}
}

func (p *PrivateHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// here user is authenticated and we should write
	ctx := request.Context()
	userId := ctx.Value("user_id")

	fmt.Println("user id read from context", userId)
	return
	// here we should start to consume user topic
	privateMessagesChan := make(chan kafkaManager.KafkaMessage)
	go kafkaManager.Consume(ctx, configs.WebSocketPrivateTopic, privateMessagesChan)

	// todo: remove this
	time.Sleep(time.Second * 2)

	// upgrade this connection to a WebSocket
	webSocket, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Println(err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	//userId := "10"
	// user id should be an uuid

	lst, ok := privateConnections.Load(userId)
	if ok {
		_slice := lst.([]*websocket.Conn)
		privateConnections.Store(userId, append(_slice, webSocket))
	} else {
		privateConnections.Store(userId, []*websocket.Conn{webSocket})
	}

}

func (p *PrivateHandler) GetChannelType() string {
	return "private"
}

func (p *PrivateHandler) WriteMessage(message string) {
	privateConnections.Range(func(key, value interface{}) bool {
		for _, conn := range value.([]*websocket.Conn) {
			p.writeToSocket(conn, message)
		}
		return true
	})
}

func (p *PrivateHandler) writeToSocket(socketConn *websocket.Conn, msg string) {
	messageByte := []byte(msg)
	messageType := 1
	if err := socketConn.WriteMessage(messageType, messageByte); err != nil {
		log.Println(err)
		return
	}
}
