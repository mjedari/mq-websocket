package kafkaManager

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"time"
	"websocket/configs"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
)

var Host string
var Port string

type KafkaMessage struct {
	Value         string
	CorrelationId string
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
	log.Println("consuming topic : ", topic, responseChan)
	sp := Host + ":" + Port
	readerconfig := kafka.ReaderConfig{
		Brokers:       []string{sp},
		Topic:         topic,
		GroupID:       configs.WebSocketKafkaGroup,
		RetentionTime: time.Second,
	}

	r := kafka.NewReader(readerconfig)
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("message err", err)
			break
		}
		headers, _ := json.Marshal(msg.Headers)
		log.Println("message received", topic, string(msg.Value), string(headers), getCorrelationId(msg.Headers))
		responseChan <- KafkaMessage{Value: string(msg.Value), CorrelationId: string(getCorrelationId(msg.Headers))}
	}
}

func getCorrelationId(headers []protocol.Header) string {
	for _, v := range headers {
		if v.Key == configs.CORRELATION_ID_KEY {
			return string(v.Value)
		}
	}
	return ""
}

func Produce(ctx context.Context, topic string, key string, requestId string, message string, responseTopic string) {
	w := &kafka.Writer{
		Addr:        kafka.TCP(Host + ":" + Port),
		Topic:       topic,
		MaxAttempts: 1,
		Balancer:    &kafka.LeastBytes{},
		BatchSize:   1,
	}
	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(message),
			Headers: []protocol.Header{
				{Key: configs.CORRELATION_ID_KEY, Value: []byte(requestId)},
				{Key: configs.WEBSOCKET_RESPONSE_KEY, Value: []byte(responseTopic)},
			},
		},
	)
	if err != nil {
		log.Println("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Println("failed to close", err)
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
