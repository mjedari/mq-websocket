package authenticationService

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
	"websocket/configs"
	"websocket/kafkaManager"
)

var channels sync.Map

type Headers struct {
	AccessToken string      `json:"access_token"`
	Headers     http.Header `json:"headers"`
}

type authenticationServerResponse struct {
	Data authenticationData `json:"data"`
}

type authenticationData struct {
	UserId   string `json:"user_id"`
	DeviceId string `json:"device_id"`
}

func init() {
	go consumeAuthenticationTopic()
}

func Authenticate(headers http.Header, requestId string, ctx context.Context) (string, string, error) {
	requestChannel := make(chan string)
	channels.Store(requestId, requestChannel)
	jsonValue, err := json.Marshal(
		Headers{AccessToken: headers.Get("Authorization"), Headers: headers},
	)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	go kafkaManager.Produce(
		ctx, configs.AUTHENTICATION_TOPIC,
		configs.AUTHENTICATION_KEY, requestId, string(jsonValue),
		configs.WEBSOCKET_AUTHENTICATION_TOPIC,
	)
	select {
	case resp := <-requestChannel:
		{
			var respModel authenticationServerResponse
			json.Unmarshal([]byte(resp), &respModel)
			return respModel.Data.UserId, respModel.Data.DeviceId, nil
		}
	case <-time.After(configs.AuthTimeout * time.Second):
		{
			log.Println("Authentication server is down !!")
			return "", "", errors.New("authentication timeout")
		}
	}
}

func consumeAuthenticationTopic() {
	userIdChannel := make(chan kafkaManager.KafkaMessage)
	go kafkaManager.Consume(context.Background(), configs.WEBSOCKET_AUTHENTICATION_TOPIC, configs.AuthenticationKafkaGroup, userIdChannel)
	for {
		kafkaMessage := <-userIdChannel
		c, ok := channels.Load(kafkaMessage.CorrelationId)
		if ok {
			channel := c.(chan string)
			channel <- kafkaMessage.Value
		} else {
			log.Println(c)
		}
	}
}
