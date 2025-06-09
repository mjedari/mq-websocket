package clients

import "github.com/mjedari/mq-websocket/domain/contracts"

type PublicClient struct {
	BaseClient
}

func (p PublicClient) Check(s string) bool {
	return true
}

func NewPublicClient(socket contracts.ISocket) *PublicClient {
	baseClient := NewBaseClient(socket)
	return &PublicClient{*baseClient}
}
