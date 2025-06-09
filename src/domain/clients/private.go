package clients

import (
	"github.com/mjedari/mq-websocket/domain/contracts"
)

type PrivateClient struct {
	BaseClient
	UserId   string
	DeviceId string
}

func NewPrivateClient(socket contracts.ISocket, userId string, deviceId string) *PrivateClient {
	baseClient := NewBaseClient(socket)
	return &PrivateClient{BaseClient: *baseClient, UserId: userId, DeviceId: deviceId}
}

func (c PrivateClient) Check(userId string) bool {
	return c.UserId == userId
}
