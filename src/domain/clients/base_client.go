package clients

import (
	"context"
	"github.com/google/uuid"
	"log"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
)

type BaseClient struct {
	Id     uuid.UUID
	Socket contracts.ISocket
	Send   chan []byte
	Close  chan bool
}

func (c *BaseClient) Check(s string) bool {
	//TODO implement me
	panic("implement me")
}

func (c *BaseClient) RemoveConnection() {
	c.Socket = nil
}

func (c *BaseClient) GetId() uuid.UUID {
	return c.Id
}

func NewBaseClient(socket contracts.ISocket) *BaseClient {
	id, _ := uuid.NewUUID()
	return &BaseClient{Id: id, Socket: socket, Send: make(chan []byte), Close: make(chan bool)}
}

func (c *BaseClient) WriteOnConnection(ctx context.Context) {
	for {
		select {
		case message, ok := <-c.Send:
			if c.Socket == nil {
				return
			}
			if !ok {
				c.Socket.WriteMessage(contracts.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(contracts.TextMessage, message)
		case <-c.Close:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (c *BaseClient) ReadFromClient(ctx context.Context) {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		default:
			_, _, err := c.Socket.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
		case <-ctx.Done():
			return
		}
	}

}

func (c *BaseClient) Leave() {
	c.Close <- true
}

func (c *BaseClient) SendMessage(message []byte) {
	c.Send <- message
}
