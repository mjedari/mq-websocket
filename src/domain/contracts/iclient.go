package contracts

import (
	"context"
	"github.com/google/uuid"
)

type IClient interface {
	WriteOnConnection(ctx context.Context)
	RemoveConnection()
	ReadFromClient(ctx context.Context)
	Leave()
	Check(string) bool
	GetId() uuid.UUID
	SendMessage([]byte)
}

type ISocket interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	Close() error
}

const (
	TextMessage   = 1
	BinaryMessage = 2
	CloseMessage  = 8
	PingMessage   = 9
	PongMessage   = 10
)
