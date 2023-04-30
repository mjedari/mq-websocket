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
	"websocket/refactor/hub"
)

const PollingTimeout = 100  // unit: ms
const NumberOfConsumers = 1 // unit: ms

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

func Consume(ctx context.Context, topic string, responseChan chan KafkaMessage, privateChan chan hub.PrivateMessage) {
	logrus.Infof("consuming topic %s: %v \n", topic, responseChan)

	for {
		run := true
		consumer, err := createNewConsumer()
		if err != nil {
			log.Println("error in consuming topic", topic, err)
			panic(err)
		}

		if subscribeErr := consumer.SubscribeTopics([]string{topic}, nil); subscribeErr != nil {
			logrus.Error("failed to close subscriber:", subscribeErr)
		}

		for run == true {
			select {
			case <-ctx.Done():
				return

			default:
				ev := consumer.Poll(PollingTimeout)
				if ev == nil {
					continue
				}

				switch e := ev.(type) {
				case *kafka.Message:
					message := KafkaMessage{
						Value:         string(e.Value),
						CorrelationId: getCorrelationId(e.Headers),
					}

					userId := getUserId(e.Headers)
					if userId != nil {
						fmt.Println("this is private message")
						msg := hub.PrivateMessage{
							UserId:  getUserId(e.Headers),
							Room:    e.Key,
							Message: e.Value,
						}

						privateChan <- msg
						continue
					}

					headers, _ := json.Marshal(e.Headers)
					log.Println("message received", topic, string(e.Value), string(headers), getCorrelationId(e.Headers))
					responseChan <- message

				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)

					if e.Code() == kafka.ErrMaxPollExceeded {
						run = false
						consumer.Close()
						fmt.Printf("closeing consumer...")
					}

				default:
					fmt.Printf("Ignored %v\n", e)
				}
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

func getUserId(headers []kafka.Header) []byte {
	for _, v := range headers {
		if v.Key == "user-id" {
			return v.Value
		}
	}
	return nil
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
		"auto.offset.reset":     "latest",
		"broker.address.family": "v4",
		"max.poll.interval.ms":  600000,
	})

	if err != nil {
		return nil, err
	}

	fmt.Printf("create new consumer: %v\n", consumer)
	return consumer, nil
}

func generateMessageHeaders(requestId string, responseTopic string) []kafka.Header {
	return []kafka.Header{
		{Key: configs.CORRELATION_ID_KEY, Value: []byte(requestId)},
		{Key: configs.WEBSOCKET_RESPONSE_KEY, Value: []byte(responseTopic)},
	}
}
