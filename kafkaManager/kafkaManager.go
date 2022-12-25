package kafkaManager

import (
	"context"
	"log"
	"net"
	"os"
	"time"
	"websocket/configs"

	"github.com/segmentio/kafka-go"
)

var Host string
var Port string

type KafkaMessage struct {
	Value string
}

func init() {
	SetKafkaServerAddress()
}

func CreateTopic(ctx context.Context, topic string) {
	var controllerConn *kafka.Conn
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(Host, Port))
	if err != nil {
		log.Println(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Println(err.Error())
	}
}

func Consume(ctx context.Context, topic string, group string, responseChan chan KafkaMessage) {
	sp := Host + ":" + Port
	readerconfig := kafka.ReaderConfig{
		Brokers:        []string{sp},
		Topic:          topic,
		Partition:      0,
		CommitInterval: time.Microsecond * 10,
		GroupID:        configs.WebSocketKafkaGroup,
	}

	r := kafka.NewReader(readerconfig)
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			continue
		}
		log.Println("message received", topic, string(msg.Value))
		responseChan <- KafkaMessage{Value: string(msg.Value)}
		go r.CommitMessages(ctx, msg)
	}
}

func SetKafkaServerAddress() {
	Host = os.Getenv("KAFKA_HOST")
	Port = os.Getenv("KAFKA_PORT")
	if Host == "" {
		Host = "localhost"
		log.Println("Using default kafka host")
	}
	if Port == "" {
		Port = "9092"
		log.Println("Using default kafka port")
	}
}
