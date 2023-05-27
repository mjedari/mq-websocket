package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	uid, _ := uuid.NewUUID()
	userId := r.Context().Value("user_id").(string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("user if from context", userId)
	conn, socketErr := upgrader.Upgrade(w, r, nil)
	if socketErr != nil {
		log.Println(socketErr)
		return
	}

	newClient := rooms.NewClient(uid, userId, conn)

	go newClient.WriteOnConnection(ctx)

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
				clientRoom, _ := h.hub.GetClientRoom(client.Id)
				if clientRoom != nil {
					clientRoom.Leave(client)
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
			// h.hub.GetClientRoom(clientId)
			err = h.hub.SetClientRoom(client.Id, r)
			if err != nil {
				// log
			}

		case "unsubscribe":
			clientRoom, _ := h.hub.GetClientRoom(client.Id)
			if clientRoom != nil {
				fmt.Println("unsubscribed to channel:", clientRoom)
				clientRoom.GetClients().Delete(client)
				_ = h.hub.RemoveClientRoom(client.Id)
			}

		case "publish":
			//clientRoom, _ := h.hub.GetClientRoom(client.Id)
			//if clientRoom != nil {
			//	clientRoom.Broadcast([]byte(msg.Data))
			//
			//}
			break
		}
	}
	return nil
}
