package wsHandler

import (
	"context"
	"github.com/google/uuid"
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

		ctx := r.Context()
		correlationId := uuid.New().String()
		userId, deviceId, err := authenticationService.Authenticate(r.Header, correlationId, ctx)
		if err != nil {
			return
		}

		context.WithValue(ctx, "user_id", userId)
		context.WithValue(ctx, "device_id", deviceId)

		// authenticate here
		next.ServeHTTP(w, r)
	})
}
