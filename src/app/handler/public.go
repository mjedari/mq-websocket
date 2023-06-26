package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/clients"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/rooms"
)

const PublicRoom = "public"

var publicUpgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type PublicHandler struct {
	hub        *hub.Hub
	monitoring contracts.IMonitoring
}

func NewPublicHandler(hub *hub.Hub, monitoring contracts.IMonitoring) *PublicHandler {
	return &PublicHandler{hub: hub, monitoring: monitoring}
}

func (h PublicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, socketErr := publicUpgrader.Upgrade(w, r, nil)
	if socketErr != nil {
		log.Println(socketErr)
		return
	}

	newClient := clients.NewPublicClient(conn)
	// h.tracer.addNewClient()
	defer func() {
		conn.Close()
	}()

	go newClient.WriteOnConnection(ctx)

	if err := h.Handle(ctx, newClient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (h PublicHandler) Handle(ctx context.Context, client *clients.PublicClient) error {
	r, err := h.hub.GetRoom(PublicRoom, func(name string) (contracts.IRoom, error) {
		return rooms.NewRoom(name)
	})

	if err != nil {
		return err
	}

	h.subscribeToRoom(client, r)

	defer func() {
		h.unSubscribeFromRoom(client, r)
		if client.Socket != nil {
			client.Socket.Close()
		}
	}()

	// wait for client if it wants to close connection
	for {
		select {
		default:
			_, _, readErr := client.Socket.ReadMessage()
			if readErr != nil {
				fmt.Println("receive message from client: ", client.GetId(), readErr)
				// we can decide what to do *additionally* with every close error be received by this function
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (h PublicHandler) subscribeToRoom(client *clients.PublicClient, room contracts.IRoom) {
	fmt.Println("subscribed to channel:", room.GetName())
	h.monitoring.AddClientToRoom(room.GetName())
	room.SetClient(client)
	h.hub.SetClientRoom(client.GetId(), room)
}

func (h PublicHandler) unSubscribeFromRoom(client *clients.PublicClient, room contracts.IRoom) {
	fmt.Println("unsubscribed from channel:", room.GetName())
	h.monitoring.RemoveClientFromRoom(room.GetName())
	client.RemoveConnection()
	h.hub.RemoveClientRoom(client.GetId())
	room.Leave(client)
}
