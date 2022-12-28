package wsHandler

import (
	"context"
	"log"
	"net/http"
	"websocket/authenticationService"
	"websocket/configs"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  configs.ReadBufferSize,
	WriteBufferSize: configs.WriteBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var connections []*websocket.Conn

func WsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request :", r.Method, r.URL.Path)

	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	ctx := context.Background()
	corolationId := "12" // TODO: set corolationId per request on websocket
	userId, deviceId, err := authenticationService.Authenticate(r.Header, corolationId, ctx)
	if err != nil {
		CloseConnection(ws)
		return
	}
	log.Println("userId: ", userId)
	log.Println("deviceId: ", deviceId)

	connections = append(connections, ws)
	go reader(ws)
}

func PublicWriter(msg string) {
	for _, conn := range connections {
		writer(conn, msg)
	}
}

func writer(conn *websocket.Conn, msg string) {
	messageByte := []byte(msg)
	messageType := 1
	if err := conn.WriteMessage(messageType, messageByte); err != nil {
		log.Println(err)
		return
	}
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			CloseConnection(conn)
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))
	}
}

func CloseConnection(conn *websocket.Conn) {
	defer conn.Close()
	for i, connItem := range connections {
		if connItem == conn {
			sz := len(connections)
			connections[i] = connections[sz-1]
			connections = connections[:sz-1]
			return
		}
	}
}
