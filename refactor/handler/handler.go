package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"websocket/refactor/hub"
	"websocket/refactor/room"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

type ChannelHandler struct {
	hub *hub.Hub
}

func NewChannelHandler(hub *hub.Hub) *ChannelHandler {
	return &ChannelHandler{hub: hub}
}

func (h ChannelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got Right Place")
	userId := r.Context().Value("user_id").(string)

	fmt.Println("user if from context", userId)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	newClient := room.NewClient(userId, conn)

	go newClient.WriteOnConnection()

	h.Handle(conn, newClient)

}

func (h ChannelHandler) Handle(conn *websocket.Conn, client *room.Client) {
	for {
		_, p, readErr := conn.ReadMessage()
		if readErr != nil {
			fmt.Println("receive message from client", client, readErr)

			if websocket.IsCloseError(readErr, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if client.Room != nil {
					client.Room.Leave(client)
				}
			}
			break
		}

		var msg Message
		if jsonErr := json.Unmarshal(p, &msg); jsonErr != nil {
			continue
		}

		switch msg.Action {
		case "subscribe":
			// can send a filter to remove sensitive information
			r := h.hub.GetRoom(msg.Channel, nil)
			fmt.Println("subscribed to channel:", r)
			r.GetClients()[client] = true
			client.Room = r

		case "unsubscribe":
			if client.Room != nil {
				fmt.Println("unsubscribed to channel:", client.Room)
				delete(client.Room.GetClients(), client)
				client.Room = nil
			}
		case "publish":
			if client.Room != nil {
				client.Room.Broadcast([]byte(msg.Data))
			}
		}
	}
}
