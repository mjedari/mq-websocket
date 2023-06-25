package messaging

type ResponseFunction func() ([]byte, error)

type ProduceMessage struct {
	Topic         string
	Key           []byte
	RequestId     string
	Message       []byte
	ResponseTopic string
}

func (p ProduceMessage) GetTopic() string {
	return p.Topic
}

func (p ProduceMessage) GetKey() []byte {
	return p.Key
}

func (p ProduceMessage) GetMessage() []byte {
	return p.Message
}

func (p ProduceMessage) GetResponseTopic() string {
	return p.ResponseTopic
}

func (p ProduceMessage) GetRequestId() string {
	return p.RequestId
}

type Message struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}
