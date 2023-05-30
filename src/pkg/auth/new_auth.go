package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"repo.abanicon.com/abantheter-microservices/websocket/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/broker"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/contracts"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/hub"
	"time"
)

type Auth struct {
	storage  contracts.IStorage
	receiver chan broker.ResponseMessage
	kafka    contracts.IBroker
	hub      *hub.Hub
}

func NewAuth(storage contracts.IStorage, kafka contracts.IBroker, hub *hub.Hub) *Auth {
	return &Auth{storage: storage, receiver: make(chan broker.ResponseMessage), kafka: kafka, hub: hub}
}

type Receiver struct {
	Action  string
	Message []byte
}

type UserToken struct {
	ID          string `json:"user_id"`
	DeviceID    string `json:"device_id"`
	AccessToken string `json:"access_token"`
}

type AuthResponse struct {
	Token UserToken
}

type AuthLogoutResponse struct {
	UserId string `json:"user_id"`
	Token  string `json:"access_token"`
}

type AccessToken struct {
	token string
}

func NewAccessToken(token string) (*AccessToken, error) {
	accessToken := AccessToken{
		token: token,
	}

	if err := accessToken.validate(); err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (accessToken AccessToken) validate() error {
	if accessToken.token == "" {
		return errors.New("invalid token")
	}

	// todo: validate all
	// validate jwt
	return nil
}

func (a *Auth) Authenticate(ctx context.Context, token string) (*UserToken, error) {
	accessToken, err := NewAccessToken(token)
	if err != nil {
		return nil, err
	}

	value := a.storage.Fetch(ctx, accessToken.token)
	if value == nil {
		return nil, errors.New("user not found")
	}

	var user UserToken
	if err := json.Unmarshal(value, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *Auth) Consume(ctx context.Context) {
	// TODO: need to generalize it
	defer fmt.Println("closing auth listener ...")
	topic := configs.Config.AuthServer.AuthenticationPrivateTopic

	go a.kafka.ConsumeAuth(ctx, topic, a.receiver)

	// waits for receiver
	for {
		select {
		case response := <-a.receiver:
			a.handleResponse(ctx, response)
		case <-ctx.Done():
			return
		}
	}
}

func (a *Auth) handleResponse(ctx context.Context, response broker.ResponseMessage) {

	switch response.Key {
	case configs.Config.AuthServer.LoginKey:
		var authResponse AuthResponse
		_ = json.Unmarshal(response.Message, &authResponse)

		fmt.Println("logging in ...", authResponse.Token.ID)
		payload, err := json.Marshal(authResponse.Token)
		if err != nil {
			//
		}

		if err := a.storage.Store(ctx, authResponse.Token.AccessToken, string(payload), time.Hour); err != nil {
			// handle err
		}

	case configs.Config.AuthServer.LogoutKey:
		// logout process
		var authLogoutResponse AuthLogoutResponse
		_ = json.Unmarshal(response.Message, &authLogoutResponse)
		fmt.Println("logging out ...", authLogoutResponse.UserId)

		// TODO: unsubscribe form channel
		//send unsubscbe signal to user id
		a.hub.LeaveClient(ctx, authLogoutResponse.UserId)

		if err := a.storage.Delete(ctx, authLogoutResponse.Token); err != nil {
			// handle err
			logrus.Error(err)
		}
	}
}
