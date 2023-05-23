package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/rooms"
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

type PrivateHandler struct {
	hub *hub.Hub
}

func NewPrivateHandler(hub *hub.Hub) *PrivateHandler {
	return &PrivateHandler{hub: hub}
}

func (h PrivateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)

	fmt.Println("user if from context", userId)
	conn, socketErr := upgrader.Upgrade(w, r, nil)
	if socketErr != nil {
		log.Println(socketErr)
		return
	}

	newClient := rooms.NewClient(userId, conn)

	go newClient.WriteOnConnection()

	if err := h.Handle(conn, newClient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h PrivateHandler) Handle(conn *websocket.Conn, client *rooms.Client) error {
	for {
		_, p, readErr := conn.ReadMessage()
		if readErr != nil {
			fmt.Println("receive message from client", client, readErr)

			if websocket.IsCloseError(readErr, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
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
			// can initiate a filtered room to remove sensitive information
			r, err := h.hub.GetRoom(msg.Channel, func(name string) (rooms.IRoom, error) {
				return rooms.NewRoom(name)
			})

			if err != nil {
				return err
			}

			fmt.Println("subscribed to channel:", r.GetName())
			r.GetClients().Store(client, true)

			client.Room = r
		case "unsubscribe":
			if client.Room != nil {
				fmt.Println("unsubscribed to channel:", client.Room)
				client.Room.GetClients().Delete(client)
				client.Room = nil
			}
		case "publish":
			if client.Room != nil {
				client.Room.Broadcast([]byte(msg.Data))
			}
		}
	}
	return nil
}
