package broker

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"
	"repo.abanicon.com/abantheter-microservices/websocket/app/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/hub"
)

const PollingTimeout = 100 // unit: ms

type Kafka struct {
	*kafka.AdminClient
}

func NewKafka(config configs.KafkaConfig) (*Kafka, error) {
	client, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers":     config.Host + ":" + config.Port,
		"broker.address.family": "v4",
		//"debug":                 "broker,protocol,feature",
	})
	if err != nil {
		return nil, err
	}

	return &Kafka{AdminClient: client}, nil
}

type ResponseMessage struct {
	Key     string
	Message []byte
}

type ProduceMessage struct {
	Topic         string
	Key           string
	RequestId     string
	Message       string
	ResponseTopic string
}

func (k *Kafka) CreateTopics(ctx context.Context, topics []string, partitions, replicationFactor int) error {
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

func (k *Kafka) Consume(ctx context.Context, topic string, publicResponseFunction, privateResponseFunction func(header, key, value []byte)) {
	defer fmt.Println("Closing kafka consumer...")
	logrus.Infof("consuming topic %s: \n", topic)
	groupId := configs.Config.Kafka.Group

	for {
		run := true
		consumer, err := createNewConsumer(groupId)
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
					//message := hub.PublicMessage{
					//	Value:         string(e.Value),
					//	CorrelationId: getCorrelationId(e.Headers),
					//}

					userId := getUserId(e.Headers)
					if userId != nil {
						fmt.Println("this is private message")
						msg := hub.PrivateMessage{
							UserId:  getUserId(e.Headers),
							Room:    e.Key,
							Message: e.Value,
						}

						// have an function with channel inside
						privateResponseFunction(msg.UserId, msg.Room, msg.Message)
						//privateChan <- msg
						continue
					}

					//headers, _ := json.Marshal(e.Headers)
					// todo: print in just debug mode
					//log.Println("message received", topic, string(e.Value), string(headers), getCorrelationId(e.Headers))
					publicResponseFunction(nil, e.Key, e.Value)
					//responseChan <- message

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

func (k *Kafka) Produce(ctx context.Context, message contracts.IBrokerMessage) {
	producer, pErr := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": configs.Config.Kafka.Host + ":" + configs.Config.Kafka.Port})
	if pErr != nil {
		panic(pErr)
	}

	topic := message.GetTopic()
	newMessage := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            message.GetKey(),
		Value:          message.GetMessage(),
		Headers:        generateMessageHeaders(message.GetRequestId(), message.GetResponseTopic()),
	}

	if err := producer.Produce(&newMessage, nil); err != nil {
		log.Print("failed to write messages:", err)
	}

	defer producer.Close()
	producer.Flush(15 * 1000)
}

func createNewConsumer(groupId string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": net.JoinHostPort(configs.Config.Kafka.Host, configs.Config.Kafka.Port),
		"group.id":          groupId,
		"auto.offset.reset": "latest",
		//"broker.address.family": "v4",
		//"max.poll.interval.ms":  600000,
	})

	if err != nil {
		return nil, err
	}

	fmt.Printf("create new consumer: %v\n", consumer)
	return consumer, nil
}

func createHealthConsumer() (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": net.JoinHostPort(configs.Config.Kafka.Host, configs.Config.Kafka.Port),
		"group.id":          configs.Config.Kafka.Group,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	fmt.Printf("create health consumer: %v\n", consumer)
	return consumer, nil
}

func generateMessageHeaders(requestId string, responseTopic string) []kafka.Header {
	return []kafka.Header{
		{Key: configs.Config.Kafka.CorrelationIdKey, Value: []byte(requestId)},
		{Key: configs.Config.Kafka.ResponseTopic, Value: []byte(responseTopic)},
	}
}

func (k *Kafka) ConsumeHealth(ctx context.Context, topic string) ([]byte, error) {
	newConsumer, err := createHealthConsumer()
	if err != nil {
		log.Println("error in consuming topic", topic, err)
		return nil, err
	}

	defer func() {
		if closeErr := newConsumer.Close(); closeErr != nil {
			logrus.Error("failed to close reader:", closeErr)
		}
	}()

	if consumeErr := newConsumer.SubscribeTopics([]string{topic}, nil); consumeErr != nil {
		return nil, err
	}

	for {
		select {
		default:
			msg, consumeErr := newConsumer.ReadMessage(-1)
			if consumeErr != nil {
				logrus.Errorf("read message error on topic %s: %v\n", topic, consumeErr)
				return nil, consumeErr
			}
			log.Println("message received", topic, string(msg.Value))
			return msg.Value, nil
		case <-ctx.Done():
			return nil, nil
		}
	}
}

func (k *Kafka) ConsumeAuth(ctx context.Context, topic string, authResponseFunction func(header, key, value []byte)) {
	defer fmt.Println("Closing kafka auth consumer...")
	logrus.Infof("consuming topic %s: %v \n", topic, authResponseFunction)
	groupId := "auth-consumer-group"
	for {
		run := true
		consumer, err := createNewConsumer(groupId)
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
					log.Println("message received", topic, string(e.Value))
					authResponseFunction(nil, e.Key, e.Value)
					// TODO: handle logout keys also

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
