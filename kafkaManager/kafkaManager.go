package kafkaManager

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"websocket/configs"
)

var Host string
var Port string

type KafkaMessage struct {
	Value         string
	CorrelationId string
}

type ProduceMessage struct {
	Topic         string
	Key           string
	RequestId     string
	Message       string
	ResponseTopic string
}

func init() {
	SetKafkaServerAddress()
}

func CreateTopic(ctx context.Context, topics []string, partitions, replicationFactor int) error {
	var topicConfigs []kafka.TopicSpecification

	for _, topic := range topics {
		t := kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replicationFactor,
		}
		topicConfigs = append(topicConfigs, t)
	}

	newAdmin, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers":     Host + ":" + Port,
		"broker.address.family": "v4",
		"debug":                 "broker,protocol,feature",
	})
	defer newAdmin.Close()

	result, err := newAdmin.CreateTopics(ctx, topicConfigs)
	if err != nil && len(result) == 0 {
		log.Println(err)
		panic(fmt.Sprintf("error in topic creation: %s", err.Error()))
	}

	logrus.Info("topics were created: ", topics)
	return nil
}

func Consume(ctx context.Context, topic string, responseChan chan KafkaMessage) {
	logrus.Infof("consuming topic %s: %v \n", topic, responseChan)

	consumer, err := createNewConsumer()
	if err != nil {
		log.Println("error in consuming topic", topic, err)
		panic(err)
	}

	defer func() {
		logrus.Info("closing consumer...\n")
		if closeErr := consumer.Close(); closeErr != nil {
			logrus.Error("failed to close consumer:", closeErr)
		}
	}()

	if subscribeErr := consumer.SubscribeTopics([]string{topic}, nil); subscribeErr != nil {
		// subscribeErr
		logrus.Error("failed to close subscriber:", subscribeErr)
	}

	run := true
	// ToDo: change the code https://github.com/confluentinc/confluent-kafka-go/blob/master/examples/consumer_example/consumer_example.go
	for run == true {
		select {
		case <-ctx.Done():
			return

		default:
			//// ToDO: use switch case to handle err and messages
			//msg, readErr := consumer.ReadMessage(time.Second)
			//if readErr != nil {
			//	logrus.Errorf("read message error on topic %s: %v\n", topic, err)
			//	continue
			//}
			//message := KafkaMessage{
			//	Value:         string(msg.Value),
			//	CorrelationId: getCorrelationId(msg.Headers),
			//}
			//
			//headers, _ := json.Marshal(msg.Headers)
			//log.Println("message received", topic, string(msg.Value), string(headers), getCorrelationId(msg.Headers))
			//responseChan <- message
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				message := KafkaMessage{
					Value:         string(e.Value),
					CorrelationId: getCorrelationId(e.Headers),
				}
				headers, _ := json.Marshal(e.Headers)
				log.Println("message received", topic, string(e.Value), string(headers), getCorrelationId(e.Headers))
				responseChan <- message

			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}

func getCorrelationId(headers []kafka.Header) string {
	for _, v := range headers {
		if v.Key == configs.CORRELATION_ID_KEY {
			return string(v.Value)
		}
	}
	return ""
}

func Produce(ctx context.Context, message ProduceMessage) {
	producer, pErr := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": Host + ":" + Port})
	if pErr != nil {
		panic(pErr)
	}

	newMessage := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &message.Topic, Partition: kafka.PartitionAny},
		Key:            []byte(message.Key),
		Value:          []byte(message.Message),
		Headers:        generateMessageHeaders(message.RequestId, message.ResponseTopic),
	}

	if err := producer.Produce(&newMessage, nil); err != nil {
		log.Print("failed to write messages:", err)
	}

	defer producer.Close()

	producer.Flush(15 * 1000)
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

func createNewConsumer() (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     Host + ":" + Port,
		"group.id":              configs.WebSocketKafkaGroup,
		"auto.offset.reset":     "earliest",
		"broker.address.family": "v4",
	})

	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func generateMessageHeaders(requestId string, responseTopic string) []kafka.Header {
	return []kafka.Header{
		{Key: configs.CORRELATION_ID_KEY, Value: []byte(requestId)},
		{Key: configs.WEBSOCKET_RESPONSE_KEY, Value: []byte(responseTopic)},
	}
}
