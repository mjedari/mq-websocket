package wsHandler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
	"websocket/configs"
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
	//return

	//go func() {
	//	// here we should start to consume user topic
	//	privateMessagesChan := make(chan kafkaManager.KafkaMessage)
	//	go kafkaManager.Consume(ctx, p.GetBrokerTopic(), privateMessagesChan)
	//
	//	for {
	//		resp := <-privateMessagesChan
	//		p.WriteMessage(resp.Value)
	//	}
	//}()

	go func() {
		for {
			<-time.Tick(time.Second * 2)
			p.WriteMessage("Hi this is message")
		}
	}()

	// todo: remove this
	time.Sleep(time.Second * 2)

	// upgrade this connection to a WebSocket
	webSocket, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Println(err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	//go reader(webSocket)

	fmt.Println("request upgrated")
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

func (p *PrivateHandler) GetBrokerTopic() string {
	return configs.WebSocketPublicTopic
}

func (p *PrivateHandler) WriteMessage(message string) {
	privateConnections.Range(func(key, value interface{}) bool {
		for _, conn := range value.([]*websocket.Conn) {
			p.writeToSocket(conn, "1", message)
		}
		return true
	})
}

func (p *PrivateHandler) writeToSocket(socketConn *websocket.Conn, userId, message string) {
	messageByte := []byte(message)
	if err := socketConn.WriteMessage(websocket.TextMessage, messageByte); err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("message sent to user: %v \n", userId)
}
