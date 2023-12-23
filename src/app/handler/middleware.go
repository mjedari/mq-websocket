package handler

import (
	"context"
	"errors"
	"net"
	"net/http"
	"repo.abanicon.com/abantheter-microservices/websocket/app/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/app/wiring"
	"repo.abanicon.com/public-library/glogger"
)

func SocketValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request has the correct headers for a WebSocket connection
		//if r.Header.Get("Upgrade") != "websocket" || r.Header.Get("Connection") != "Upgrade" {
		//	http.Error(w, "Not a WebSocket request", http.StatusBadRequest)
		//	return
		//}

		// check the origin of the request
		if err := checkRequestOrigin(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: can add a token as a validation to make us sure that the request comes from our servers
		next.ServeHTTP(w, r)
	})
}

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !configs.Config.RateLimiter.Active {
			next.ServeHTTP(w, r)
			return
		}

		rateLimiter := wiring.Wiring.GetRateLimiter()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "invalid ip address", http.StatusInternalServerError)
			return
		}

		if err := rateLimiter.Handle(ip); err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		glogger.Println("request :", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func PrivateChannelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//correlationId := uuid.New().String()
		ctx := r.Context()
		token := r.Header.Get("Authorization")
		authService := wiring.Wiring.GetAuthService()

		user, err := authService.Authenticate(ctx, token)
		if err != nil {
			http.Error(w, "token is unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "user_id", user.ID)
		ctx = context.WithValue(ctx, "device_id", user.DeviceID)
		newReq := r.WithContext(ctx)

		// authenticate here
		next.ServeHTTP(w, newReq)
	})
}

func checkRequestOrigin(r *http.Request) error {
	requestOrigin := r.Header.Get("Origin")
	validOrigins := configs.Config.Security.ValidOrigins
	if validOrigins == nil {
		return nil
	}

	for _, origin := range validOrigins {
		if requestOrigin == origin {
			return nil
		}
	}
	return errors.New("invalid request origin")
}
