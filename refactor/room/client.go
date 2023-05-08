package room

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	UserId string
	Conn   *websocket.Conn
	Send   chan []byte
	Room   IRoom
}

func NewClient(userId string, conn *websocket.Conn) *Client {
	return &Client{UserId: userId, Conn: conn, Send: make(chan []byte)}
}

func (c Client) WriteOnConnection() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)

		default:
			//
		}
	}

}

func (c Client) ReadFromClient() {
	for {
		// read in a message
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
	}

}
