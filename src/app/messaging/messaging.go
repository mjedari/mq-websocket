package messaging

import (
	"context"
	"fmt"
	"repo.abanicon.com/abantheter-microservices/websocket/app/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/messaging"
	"repo.abanicon.com/public-library/glogger"
)

type Messaging struct {
	broker     contracts.IBroker
	monitoring contracts.IMonitoring
	hub        *hub.Hub
	auth       *AuthService
}

func NewMessaging(broker contracts.IBroker, monitoring contracts.IMonitoring, hub *hub.Hub, auth *AuthService) *Messaging {
	return &Messaging{broker: broker, monitoring: monitoring, hub: hub, auth: auth}
}

func (m *Messaging) publicFunction(userId, key, value []byte) {
	msg := hub.PublicMessage{
		Room:          key,
		Message:       value,
		CorrelationId: "",
	}

	// here we have all elements of message
	// can you send it to the right channel?
	m.monitoring.MessageReceived(string(key), "public")
	m.hub.PublicReceiver <- msg
}

func (m *Messaging) privateFunction(userId, key, value []byte) {
	msg := hub.PrivateMessage{
		UserId:  userId,
		Room:    key,
		Message: value,
	}

	// here we have all elements of message
	// can you send it to the right channel?
	m.monitoring.MessageReceived(string(key), "private")
	m.hub.PrivateReceiver <- msg
}

func (m *Messaging) authFunction(userId, key, value []byte) {
	switch string(key) {
	case configs.Config.AuthServer.LoginKey:
		m.auth.login(context.Background(), value)
	case configs.Config.AuthServer.LogoutKey:
		m.auth.logout(context.Background(), value)
	}
}

func (m *Messaging) Run(ctx context.Context) {
	// create topics
	if err := m.CreateTopics(ctx); err != nil {
		panic(err)
	}

	// health check
	//if err := m.HealthCheck(ctx); err != nil {
	//	panic(err)
	//}

	// start consuming
	go m.broker.Consume(ctx, configs.Config.Topics.PublicTopic, m.publicFunction, m.privateFunction)
	go m.broker.ConsumeAuth(ctx, configs.Config.AuthServer.AuthenticationPrivateTopic, m.authFunction)

	// start handling
	m.streaming(ctx)
}

func (m *Messaging) HealthCheck(ctx context.Context) error {

	healthTopic := configs.Config.Topics.Health
	testingMessage := messaging.ProduceMessage{
		Topic:         healthTopic,
		Key:           []byte("health-check-key"),
		Message:       []byte("valid-health-check-value"),
		ResponseTopic: healthTopic,
	}

	m.broker.Produce(ctx, testingMessage)

	result, err := m.broker.ConsumeHealth(ctx, healthTopic)
	if err != nil {
		panic(err)
	}
	if string(result) == string(testingMessage.Message) {
		glogger.Info("Kafka is listening")
	}

	return nil
}

func (m *Messaging) CreateTopics(ctx context.Context) error {
	topics := []string{
		configs.Config.Topics.Health,
		configs.Config.Topics.PublicTopic,
	}
	kafkaConfig := configs.Config.Kafka

	if err := m.broker.CreateTopics(ctx, topics, kafkaConfig.Partitions, kafkaConfig.ReplicationFactor); err != nil {
		return err
	}
	return nil

}

func (m *Messaging) streaming(ctx context.Context) {
	// listening for private channel
	go m.privateStreaming(ctx)

	// listening for public channel
	go m.publicStreaming(ctx)

}

func (m *Messaging) privateStreaming(ctx context.Context) {
	defer fmt.Println("Closing private streaming ...")

	for {
		select {
		case msg := <-m.hub.PrivateReceiver:
			fmt.Println("received private message")
			fmt.Println("rooms-id ", string(msg.Room))
			fmt.Println("user-id: ", string(msg.UserId))
			fmt.Println("message: ", string(msg.Message))

			m.hub.PrivateRooms.Range(func(key, value any) bool {
				r, ok := value.(contracts.IPrivateRoom)
				if !ok {
					return false
				}
				if r.GetName() == string(msg.Room) {
					r.PrivateSend(string(msg.UserId), msg.Message)
				}

				return true
			})
		case <-ctx.Done():
			return
		}
	}
}

func (m *Messaging) publicStreaming(ctx context.Context) {
	defer fmt.Println("Closing public streaming ...")

	for {
		select {
		case msg := <-m.hub.PublicReceiver:
			// todo: this should be resilient to time.Sleep
			//time.Sleep(time.Minute)
			fmt.Println("received public message")
			fmt.Println("rooms-id ", string(msg.Room))
			fmt.Println("message: ", string(msg.Message))

			// iterate over all rooms and find which one is public through reflect package
			m.hub.PublicRooms.Range(func(key, value any) bool {
				r, ok := value.(contracts.IPublicRoom)
				if !ok {
					fmt.Println("not found public channel")
					return false
				}
				if r.GetName() == string(msg.Room) {
					r.Broadcast(msg.Message)
				}

				return true
			})

		case <-ctx.Done():
			return
		}
	}
}

func (m *Messaging) LeaveClient(ctx context.Context, userId string) {
	m.hub.PrivateRooms.Range(func(room, anyRoom any) bool {
		r, ok := anyRoom.(contracts.IRoom)
		if !ok {
			return false
		}

		r.GetClients().Range(func(cli, status any) bool {
			c := cli.(contracts.IClient)

			if c.Check(userId) {
				r.Leave(c)
			}
			return true
		})

		return true
	})
}
