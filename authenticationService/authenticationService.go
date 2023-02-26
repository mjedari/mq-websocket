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

type StandardHeader map[string]string

func Standardize(headers http.Header) StandardHeader {
	standard_headers := make(StandardHeader)
	for key := range headers {
		standard_headers[key] = headers.Get(key)
	}
	return standard_headers
}

type AuthFormat struct {
	AccessToken string   `json:"access_token"`
	MetaData    MetaData `json:"MetaData"`
}

type MetaData struct {
	Headers StandardHeader `json:"headers"`
	Ip      string         `json:"ip"`
}

type authenticationServerResponse struct {
	Data     authenticationData `json:"data"`
	MetaData authenticationMetaData
}

type authenticationMetaData struct {
	Headers map[string]any `json:"headers"`
}

type authenticationData struct {
	UserId   string `json:"user_id"`
	DeviceId string `json:"device_id"`
}

func init() {
	go consumeAuthenticationTopic()
}

func Authenticate(headers http.Header, requestId string, ctx context.Context) (string, string, error) {
	var respModel authenticationServerResponse
	requestChannel := make(chan string)
	channels.Store(requestId, requestChannel)
	jsonValue, err := json.Marshal(
		AuthFormat{
			AccessToken: headers.Get("Authorization"),
			MetaData:    MetaData{Headers: Standardize(headers), Ip: headers.Get("X-FORWARDED-FOR")},
		},
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
			log.Println("Authentication Response :", resp)
			json.Unmarshal([]byte(resp), &respModel)
			return respModel.Data.UserId, respModel.Data.DeviceId, nil
		}
	case <-time.After(configs.AuthTimeout * time.Second):
		{
			channels.Delete(requestId)
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
			log.Println("request channel not found: ", kafkaMessage.CorrelationId)
		}
	}
}
