package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/room"
)

var publicUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type PublicHandler struct {
	hub *hub.Hub
}

func NewPublicHandler(hub *hub.Hub) *PublicHandler {
	return &PublicHandler{hub: hub}
}

func (h PublicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got Right Place")

	uid, err := uuid.NewUUID()
	userId := uid.String()

	fmt.Println("user if from context", userId)
	conn, err := publicUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	newClient := room.NewClient(userId, conn)

	go newClient.WriteOnConnection()

	h.Handle(conn, newClient)

}

func (h PublicHandler) Handle(conn *websocket.Conn, client *room.Client) {
	r := h.hub.GetRoom("public", nil)
	fmt.Println("subscribed to public channel:", r)
	r.GetClients()[client] = true
	client.Room = r

	for {
		_, _, readErr := conn.ReadMessage()
		if readErr != nil {
			fmt.Println("receive message from client", client, readErr)

			if websocket.IsCloseError(readErr, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if client.Room != nil {
					client.Room.Leave(client)
				}
			}
			break
		}
	}
}
