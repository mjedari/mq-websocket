package wsHandler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"websocket/authenticationService"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		log.Println("request :", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func PrivateChannelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// authenticate here
		// => how?
		// produce a message to auth server
		//	=> wait to get it is allowed or not
		//		=> then allow to establish a connection
		//		=> or then close connection

		fmt.Println("Got authentication middleware")
		correlationId := uuid.New().String()
		ctx := r.Context()
		userId, deviceId, err := authenticationService.Authenticate(r.Header, correlationId, ctx)
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
