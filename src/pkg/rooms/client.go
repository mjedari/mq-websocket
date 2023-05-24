package rooms

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Id     uuid.UUID
	UserId string
	Conn   *websocket.Conn
	Send   chan []byte
	Close  chan bool
	//Room   IRoom
}

func NewClient(id uuid.UUID, userId string, conn *websocket.Conn) *Client {
	return &Client{Id: id, UserId: userId, Conn: conn, Send: make(chan []byte), Close: make(chan bool)}
}

func (c *Client) WriteOnConnection() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		case <-c.Close:
			return
		}
	}
}

func (c *Client) ReadFromClient() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		// read in a message
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
	}

}
