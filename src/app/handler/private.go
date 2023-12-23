package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/clients"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/messaging"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/rooms"
	"repo.abanicon.com/public-library/glogger"
	"sync"
)

var privateUpgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type PrivateHandler struct {
	hub          *hub.Hub
	monitoring   contracts.IMonitoring
	requestRoom  *contracts.IPrivateRoom
	requestRooms sync.Map
}

func NewPrivateHandler(hub *hub.Hub, monitoring contracts.IMonitoring) *PrivateHandler {
	return &PrivateHandler{hub: hub, monitoring: monitoring}
}

func (h *PrivateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	deviceId := r.Context().Value("device_id").(string)

	ctx, cancel := context.WithCancel(context.Background())

	conn, socketErr := privateUpgrader.Upgrade(w, r, nil)
	if socketErr != nil {
		glogger.Error(socketErr)
		return
	}
	defer func() {
		cancel()
		conn.Close()
	}()

	newClient := clients.NewPrivateClient(conn, userId, deviceId)

	go newClient.WriteOnConnection(ctx)

	if err := h.Handle(ctx, newClient); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h PrivateHandler) Handle(ctx context.Context, client *clients.PrivateClient) error {
	h.monitoring.AddConnection("private")

	defer func() {
		if client.Socket != nil {
			client.Socket.Close()
		}
		client.Leave()
		client.RemoveConnection()

		h.requestRooms.Range(func(key, value any) bool {
			room, ok := value.(contracts.IRoom)
			if !ok {
				return false
			}
			h.unSubscribeFromRoom(client, room)
			return true
		})

		h.monitoring.RemoveConnection("private")
	}()

	for {
		_, p, readErr := client.Socket.ReadMessage()
		if readErr != nil {
			fmt.Println("receive private connection message from client user-id: ", client.UserId, readErr)
			// todo: unsubscribe from all the rooms
			break
		}

		var msg messaging.Message
		if jsonErr := json.Unmarshal(p, &msg); jsonErr != nil {
			continue
		}

		switch msg.Action {
		case "subscribe":
			// can initiate a filtered rooms to remove sensitive information
			room, err := h.hub.GetPrivateRoom(msg.Channel, func(name string) (contracts.IPrivateRoom, error) {
				return rooms.NewPrivateRoom(name)
			})

			h.requestRooms.Store(room.GetName(), room)

			if err != nil {
				// todo: handle this: decide to return error or continue
				return err
			}
			h.subscribeToRoom(client, room)

		case "unsubscribe":
			room, err := h.hub.GetPrivateRoom(msg.Channel, nil)
			if err != nil {
				// todo: glogger the error and wait to next command
				continue
			}

			h.unSubscribeFromRoom(client, room)

		case "publish":
			// it is not featured to be implemented
		}
	}
	return nil
}

func (h *PrivateHandler) subscribeToRoom(client *clients.PrivateClient, room contracts.IRoom) {
	fmt.Println("subscribed to channel:", room.GetName())
	h.monitoring.AddClientToRoom(room.GetName())
	room.SetClient(client)
	h.hub.SetClientRoom(client.GetId(), room)
}

func (h *PrivateHandler) unSubscribeFromRoom(client *clients.PrivateClient, rooms ...contracts.IRoom) {
	for _, room := range rooms {
		fmt.Println("unsubscribed from channel:", room.GetName())
		//client.RemoveConnection()
		h.hub.RemoveClientRoom(client.GetId())
		existed := room.Leave(client)
		if existed {
			h.monitoring.RemoveClientFromRoom(room.GetName())
		}
	}
}
