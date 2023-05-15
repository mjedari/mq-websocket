package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/auth"
)

func SocketValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request has the correct headers for a WebSocket connection
		if r.Header.Get("Upgrade") != "websocket" || r.Header.Get("Connection") != "Upgrade" {
			http.Error(w, "Not a WebSocket request", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		log.Println("request :", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func PrivateChannelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got authentication middleware")
		correlationId := uuid.New().String()
		ctx := r.Context()
		userId, deviceId, err := auth.Authenticate(r.Header, correlationId, ctx)
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userId == "" {
			logrus.Error("authentication failed")
			http.Error(w, "authentication failed", http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, "user_id", userId)
		ctx = context.WithValue(ctx, "device_id", deviceId)
		newReq := r.WithContext(ctx)

		// authenticate here
		next.ServeHTTP(w, newReq)
	})
}
