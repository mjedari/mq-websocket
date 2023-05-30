package handler

import (
	"context"
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
	ReadBufferSize:  0,
	WriteBufferSize: 0, // todo: get the right size for both
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
	uid, _ := uuid.NewUUID()
	userId := uid.String()
	ctx, cancel := context.WithCancel(context.Background())

	conn, socketErr := publicUpgrader.Upgrade(w, r, nil)
	if socketErr != nil {
		log.Println(socketErr)
		return
	}

	defer func() {
		cancel()
		conn.Close()
	}()

	newClient := rooms.NewClient(uid, userId, conn)

	go newClient.WriteOnConnection(ctx)

	if err := h.Handle(ctx, conn, newClient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h PublicHandler) Handle(ctx context.Context, conn *websocket.Conn, client *rooms.Client) error {
	r, err := h.hub.GetRoom(PublicRoom, func(name string) (rooms.IRoom, error) {
		return rooms.NewRoom(name)
	})

	if err != nil {
		return err
	}

	h.subscribeToRoom(client, r)

	defer func() {
		h.unSubscribeFromRoom(client, r)
		conn.Close()
	}()

	err = h.hub.SetClientRoom(client.Id, r)
	if err != nil {
		log.Println("error", err)
	}

	// wait for client if it wants to close connection
	for {
		select {
		default:
			_, _, readErr := conn.ReadMessage()
			if readErr != nil {
				fmt.Println("receive message from client: ", client.UserId, readErr)
				// we can decide what to do *additionally* with every close error be received by this function
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (h PublicHandler) subscribeToRoom(client *rooms.Client, r rooms.IRoom) {
	fmt.Printf("subscribed to %v channel: %v\n", PublicRoom, client.UserId)
	r.GetClients().Store(client, true)
}

func (h PublicHandler) unSubscribeFromRoom(client *rooms.Client, r rooms.IRoom) {
	client.Conn = nil
	r.GetClients().Delete(client)
	clientRoom, _ := h.hub.GetClientRoom(client.Id)
	if clientRoom != nil {
		clientRoom.Leave(client)
	}
	_ = h.hub.RemoveClientRoom(client.Id)
}
