package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
	"websocket/configs"
	"websocket/kafkaManager"
	"websocket/refactor/handler"
	"websocket/refactor/hub"
	"websocket/wsHandler"

	"github.com/getsentry/sentry-go"
)

func main() {
	InitSentry()
	CreateTopics() // create required websocket topics
	privateChan := make(chan hub.PrivateMessage)

	go publicMessageManager.ReceiveMessages(privateChan)

	mux := http.NewServeMux()
	mux.HandleFunc("/", wsHandler.WsHandler)

	newHub := hub.NewHub()
	channelHandler := handler.NewChannelHandler(newHub)
	privateSocketHandler := handler.LoggerMiddleware(handler.PrivateChannelMiddleware(channelHandler))

	go newHub.Streaming(privateChan)

	mux.Handle("/private", privateSocketHandler)

	err := http.ListenAndServe(":"+configs.WebSocketPort, mux)
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
	if err := kafkaManager.CreateTopic(ctx, []string{configs.WebSocketPublicTopic}, 1, 1); err != nil {
		logrus.Error("can not create topic: ", configs.WebSocketPublicTopic)
		panic(err)
	}
}
