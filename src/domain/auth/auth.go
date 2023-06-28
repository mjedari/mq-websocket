package auth

import (
	"errors"
	"time"
)

type UserToken struct {
	ID          string    `json:"user_id"`
	DeviceID    string    `json:"device_id"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   ExpiresAt `json:"access_token_expired_at"`
}

type AuthResponse struct {
	Token UserToken
}

func (a AuthResponse) GetExpiresTime() time.Duration {
	timestamp := time.Unix(a.Token.ExpiresAt.Timestamp/1000, 0)
	return timestamp.Sub(time.Now())
}

type AuthLogoutResponse struct {
	UserId string `json:"user_id"`
	Token  string `json:"access_token"`
}

type AccessToken struct {
	Token     string
	ExpiresAt time.Duration
}

type ExpiresAt struct {
	Timestamp int64 `json:"$date"`
}

func NewAccessToken(token string) (*AccessToken, error) {
	accessToken := AccessToken{
		Token: token,
	}

	if err := accessToken.validate(); err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (accessToken AccessToken) validate() error {
	if accessToken.Token == "" {
		return errors.New("invalid token")
	}

	// todo: validate all
	// validate jwt
	return nil
}
