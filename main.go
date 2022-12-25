package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"websocket/configs"
	"websocket/kafkaManager"
	"websocket/publicMessageManager"
	"websocket/wsHandler"

	"github.com/getsentry/sentry-go"
)

func main() {
	InitSentry()
	CreateTopics() // create required websocket topics
	go publicMessageManager.ReceiveMessages()
	http.HandleFunc("/", wsHandler.WsHandler)
	err := http.ListenAndServe(":"+configs.WebSocketPort, nil)
	if err != nil {
		log.Fatal(err)
		sentry.CaptureException(err)
	}
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
