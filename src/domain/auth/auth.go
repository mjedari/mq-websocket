package auth

import (
	"errors"
)

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
	Token string
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
