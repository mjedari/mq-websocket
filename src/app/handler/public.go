package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mjedari/mq-websocket/domain/clients"
	"github.com/mjedari/mq-websocket/domain/contracts"
	"github.com/mjedari/mq-websocket/domain/hub"
	"github.com/mjedari/mq-websocket/domain/messaging"
	"github.com/mjedari/mq-websocket/domain/rooms"
	"github.com/mjedari/mq-websocket/infra/utils"
	"net/http"
	"repo.abanicon.com/public-library/glogger"
)

var publicUpgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type PublicHandler struct {
	hub          *hub.Hub
	monitoring   contracts.IMonitoring
	requestRooms *utils.SafeMap
}

func NewPublicHandler(hub *hub.Hub, monitoring contracts.IMonitoring) *PublicHandler {
	return &PublicHandler{hub: hub, requestRooms: utils.NewSafeMap(), monitoring: monitoring}
}

func (h *PublicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, socketErr := publicUpgrader.Upgrade(w, r, nil)
	if socketErr != nil {
		glogger.Error(socketErr)
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

func (h *PublicHandler) Handle(ctx context.Context, client *clients.PublicClient) error {
	h.monitoring.AddConnection("public")

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

		h.monitoring.RemoveConnection("public")
	}()

	for {
		_, p, readErr := client.Socket.ReadMessage()
		if readErr != nil {
			fmt.Println("receive public connection message: ", readErr)
			break
		}

		var msg messaging.Message
		if jsonErr := json.Unmarshal(p, &msg); jsonErr != nil {
			continue
		}

		switch msg.Action {
		case "subscribe":
			// can initiate a filtered rooms to remove sensitive information
			room, err := h.hub.GetPublicRoom(msg.Channel, func(name string) (contracts.IPublicRoom, error) {
				return rooms.NewPublicRoom(name)
			})

			// note: check this because it is not no longer a pointer *room
			h.requestRooms.Store(room.GetName(), room)
			if err != nil {
				// todo: handle this: decide to return error or continue
				return err
			}
			h.subscribeToRoom(client, room)

		case "unsubscribe":
			room, err := h.hub.GetPublicRoom(msg.Channel, nil)
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

func (h *PublicHandler) subscribeToRoom(client *clients.PublicClient, room contracts.IRoom) {
	fmt.Println("subscribed to channel:", room.GetName())
	h.monitoring.AddClientToRoom(room.GetName())
	room.SetClient(client)
	h.hub.SetClientRoom(client.GetId(), room)
}

func (h *PublicHandler) unSubscribeFromRoom(client *clients.PublicClient, rooms ...contracts.IRoom) {
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
