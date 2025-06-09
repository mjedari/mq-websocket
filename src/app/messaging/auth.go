package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mjedari/mq-websocket/domain/auth"
	"github.com/mjedari/mq-websocket/domain/contracts"
	"github.com/mjedari/mq-websocket/domain/hub"
	"repo.abanicon.com/public-library/glogger"
)

type AuthService struct {
	storage    contracts.IStorage
	monitoring contracts.IMonitoring
	hub        *hub.Hub
}

func NewAuthService(storage contracts.IStorage, monitoring contracts.IMonitoring, hub *hub.Hub) *AuthService {
	return &AuthService{storage: storage, monitoring: monitoring, hub: hub}
}

func (a *AuthService) Authenticate(ctx context.Context, token string) (*auth.UserToken, error) {
	accessToken, err := auth.NewAccessToken(token)
	if err != nil {
		a.monitoring.AuthenticationFailed()
		return nil, err
	}

	value := a.storage.Fetch(ctx, accessToken.Token)
	if value == nil {
		a.monitoring.AuthenticationFailed()
		return nil, errors.New("user not found")
	}

	var user auth.UserToken
	if decodeErr := json.Unmarshal(value, &user); decodeErr != nil {
		return nil, decodeErr
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

	if err := a.storage.Store(ctx, authResponse.Token.AccessToken, string(payload), authResponse.GetExpiresTime()); err != nil {
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
		glogger.Error(err)
	}
}
