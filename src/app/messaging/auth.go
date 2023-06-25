package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/auth"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/hub"
	"time"
)

type AuthService struct {
	storage contracts.IStorage
	hub     *hub.Hub
}

func NewAuthService(storage contracts.IStorage, hub *hub.Hub) *AuthService {
	return &AuthService{storage: storage, hub: hub}
}

func (a *AuthService) Authenticate(ctx context.Context, token string) (*auth.UserToken, error) {
	accessToken, err := auth.NewAccessToken(token)
	if err != nil {
		return nil, err
	}

	value := a.storage.Fetch(ctx, accessToken.Token)
	if value == nil {
		return nil, errors.New("user not found")
	}

	var user auth.UserToken
	if err := json.Unmarshal(value, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *AuthService) login(ctx context.Context, message []byte) {

	var authResponse auth.AuthResponse
	_ = json.Unmarshal(message, &authResponse)

	fmt.Println("logging in ...", authResponse.Token.ID)
	payload, err := json.Marshal(authResponse.Token)
	if err != nil {
		//
	}

	if err := a.storage.Store(ctx, authResponse.Token.AccessToken, string(payload), time.Hour); err != nil {
		// handle err
	}
}

func (a *AuthService) logout(ctx context.Context, message []byte) {

	// logout process
	var authLogoutResponse auth.AuthLogoutResponse
	_ = json.Unmarshal(message, &authLogoutResponse)
	fmt.Println("logging out ...", authLogoutResponse.UserId)

	// TODO: unsubscribe form channel
	//send unsubscbe signal to user id

	//a.hub.LeaveClient(ctx, authLogoutResponse.UserId) //todo: uncomment

	if err := a.storage.Delete(ctx, authLogoutResponse.Token); err != nil {
		// handle err
		logrus.Error(err)
	}
}
