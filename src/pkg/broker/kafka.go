package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"
	"repo.abanicon.com/abantheter-microservices/websocket/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/hub"
)

const PollingTimeout = 100 // unit: ms

type Kafka struct {
	*kafka.AdminClient
}

func NewKafka(config configs.KafkaConfig) (*Kafka, error) {
	client, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers":     config.Host + ":" + config.Port,
		"broker.address.family": "v4",
		"debug":                 "broker,protocol,feature",
	})
	if err != nil {
		return nil, err
	}

	return &Kafka{AdminClient: client}, nil
}

type ProduceMessage struct {
	Topic         string
	Key           string
	RequestId     string
	Message       string
	ResponseTopic string
}

func (k *Kafka) CreateTopic(ctx context.Context, topics []string, partitions, replicationFactor int) error {
	var topicConfigs []kafka.TopicSpecification
	for _, topic := range topics {
		t := kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replicationFactor,
		}
		topicConfigs = append(topicConfigs, t)
	}
	result, err := k.AdminClient.CreateTopics(ctx, topicConfigs)
	if err != nil && len(result) == 0 {
		log.Println(err)
		panic(fmt.Sprintf("error in topic creation: %s", err.Error()))
	}

	logrus.Info("topics were created: ", topics)
	return nil
}

func (k *Kafka) Consume(ctx context.Context, topic string, responseChan chan hub.KafkaMessage, privateChan chan hub.PrivateMessage) {
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
					message := hub.KafkaMessage{
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
		if v.Key == configs.Config.Kafka.CorrelationIdKey {
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
	fmt.Println("Got produce", configs.Config.Kafka.Host+":"+configs.Config.Kafka.Port)
	producer, pErr := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": configs.Config.Kafka.Host + ":" + configs.Config.Kafka.Port})
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

func createNewConsumer() (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     net.JoinHostPort(configs.Config.Kafka.Host, configs.Config.Kafka.Port),
		"group.id":              configs.Config.Kafka.Group,
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
		{Key: configs.Config.Kafka.CorrelationIdKey, Value: []byte(requestId)},
		{Key: configs.Config.Kafka.ResponseTopic, Value: []byte(responseTopic)},
	}
}

func ConsumeHealth(ctx context.Context, topic string) (string, error) {
	newConsumer, err := createNewConsumer()
	if err != nil {
		log.Println("error in consuming topic", topic, err)
		return "", err
	}

	defer func() {
		if closeErr := newConsumer.Close(); closeErr != nil {
			logrus.Error("failed to close reader:", closeErr)
		}
	}()

	if consumeErr := newConsumer.SubscribeTopics([]string{topic}, nil); consumeErr != nil {
		return "", err
	}

	for {
		msg, consumeErr := newConsumer.ReadMessage(-1)
		if consumeErr != nil {
			logrus.Errorf("read message error on topic %s: %v\n", topic, consumeErr)
			return "", consumeErr
		}
		log.Println("message received", topic, string(msg.Value))
		return string(msg.Value), nil
	}
}
