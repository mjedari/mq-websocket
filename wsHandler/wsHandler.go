package wsHandler

import (
	"context"
	"log"
	"net/http"
	"sync"
	"websocket/authenticationService"
	"websocket/configs"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  configs.ReadBufferSize,
	WriteBufferSize: configs.WriteBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var PrivateConnections sync.Map
var PublicConnections sync.Map

func WsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request :", r.Method, r.URL.Path)

	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	go reader(ws)

	ctx := context.Background()
	correlationId := uuid.New().String()
	userId, deviceId, err := authenticationService.Authenticate(r.Header, correlationId, ctx)
	if err != nil {
		CloseConnection(ws)
		return
	}
	log.Println("userId: ", userId)
	log.Println("deviceId: ", deviceId)
	if len(userId) != 0 {
		// Private Connection
		lst, ok := PrivateConnections.Load(userId)
		if ok {
			_slice := lst.([]*websocket.Conn)
			PrivateConnections.Store(userId, append(_slice, ws))
		} else {
			PrivateConnections.Store(userId, []*websocket.Conn{ws})
		}
	} else {
		// Public Connection
		lst, ok := PublicConnections.Load(userId)
		if ok {
			_slice := lst.([]*websocket.Conn)
			PublicConnections.Store(deviceId, append(_slice, ws))
		} else {
			PublicConnections.Store(deviceId, []*websocket.Conn{ws})
		}
	}
}

func PublicWriter(msg string) {
	PublicConnections.Range(func(key, value interface{}) bool {
		for _, conn := range value.([]*websocket.Conn) {
			writer(conn, msg)
		}
		return true
	})
	PrivateConnections.Range(func(key, value interface{}) bool {
		for _, conn := range value.([]*websocket.Conn) {
			writer(conn, msg)
		}
		return true
	})
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
	PublicConnections.Range(
		func(key, value interface{}) bool {
			connections := value.([]*websocket.Conn)
			for i, connItem := range connections {
				if connItem == conn {
					sz := len(connections)
					connections[i] = connections[sz-1]
					connections = connections[:sz-1]
					PublicConnections.Store(key, connections)
					return false
				}
			}
			return true
		})
	PrivateConnections.Range(
		func(key, value interface{}) bool {
			connections := value.([]*websocket.Conn)
			for i, connItem := range connections {
				if connItem == conn {
					sz := len(connections)
					connections[i] = connections[sz-1]
					connections = connections[:sz-1]
					PrivateConnections.Store(key, connections)
					return false
				}
			}
			return true
		})
}
