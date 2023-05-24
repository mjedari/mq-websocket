package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/rooms"
)

const PublicRoom = "public"

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
	fmt.Println("Got PublicHandler")
	uid, _ := uuid.NewUUID()
	userId := uid.String()

	conn, socketErr := publicUpgrader.Upgrade(w, r, nil)
	if socketErr != nil {
		log.Println(socketErr)
		return
	}

	newClient := rooms.NewClient(uid, userId, conn)

	go newClient.WriteOnConnection()

	if err := h.Handle(conn, newClient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (h PublicHandler) Handle(conn *websocket.Conn, client *rooms.Client) error {
	fmt.Println("Got PublicHandler Handle")

	r, err := h.hub.GetRoom(PublicRoom, func(name string) (rooms.IRoom, error) {
		return rooms.NewRoom(name)
	})

	if err != nil {
		return err
	}

	fmt.Printf("subscribed to %v channel: %v\n", PublicRoom, client.UserId)
	r.GetClients().Store(client, true)
	//client.Room = r
	err = h.hub.SetClientRoom(client.Id, r)
	if err != nil {
		// log
		log.Println("error", err)
	}

	// wait for client if it wants to close connection
	fmt.Println("Got PublicHandler wait for")

	for {
		_, _, readErr := conn.ReadMessage()
		if readErr != nil {
			fmt.Println("receive message from client: ", client.UserId, readErr)

			if websocket.IsCloseError(readErr, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				clientRoom, _ := h.hub.GetClientRoom(client.Id)
				if clientRoom != nil {
					clientRoom.Leave(client)
				}
			}
			break
		}
	}
	return nil
}
