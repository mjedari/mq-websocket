package main

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
	"websocket/authenticationService"
	"websocket/configs"
	"websocket/kafkaManager"
	"websocket/refactor/handler"
	"websocket/refactor/hub"
)

func main() {
	InitSentry()
	CreateTopics() // create required websocket topics

	kafkaHandler := kafkaManager.NewKafkaHandler()
	newHub := hub.NewHub(kafkaHandler)

	newHub.Streaming()
	go authenticationService.HandleAuthMessage(newHub.AuthReceiver)

	newPrivateHandler := handler.NewPrivateHandler(newHub)
	privateHandler := handler.LoggerMiddleware(handler.PrivateChannelMiddleware(newPrivateHandler))

	newPublicHandler := handler.NewPublicHandler(newHub)
	publicHandler := handler.LoggerMiddleware(newPublicHandler)

	// build endpoints
	mux := http.NewServeMux()
	mux.Handle("/", publicHandler)
	mux.Handle("/private", privateHandler)

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
