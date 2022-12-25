package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"websocket/configs"
	"websocket/kafkaManager"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  configs.ReadBufferSize,
	WriteBufferSize: configs.WriteBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	InitSentry()
	CreateTopics() // create required websocket topics
	http.HandleFunc("/ws", wsEndpoint)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func InitSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              configs.SENTRY_DSN,
		Environment:      configs.ENVIRONMENT,
		Debug:            configs.DEBUG,
		TracesSampleRate: 1.0,
		AttachStacktrace: true,
		ServerName:       "WebSocket",
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer func() {
		err := recover()

		if err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(time.Second * 5)
		}
	}()
}

func CreateTopics() {
	ctx := context.Background()
	kafkaManager.CreateTopic(ctx, configs.WebSocketPublicTopic)
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	log.Println(ws.RemoteAddr())
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	go reader(ws)
	go writer(ws)
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))
	}
}

func writer(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 2)
		messageByte := []byte("salam")
		messageType := 1
		if err := conn.WriteMessage(messageType, messageByte); err != nil {
			log.Println(err)
			return
		}
	}
}
