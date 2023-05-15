package configs

var Config Configuration

type Server struct {
	Host string
	Port string
}

type KafkaConfig struct {
	Host              string
	Port              string
	DefaultGroup      string
	Group             string
	ResponseTopic     string
	Partitions        int `mapstructure:"ResponsePartitionsCount"`
	ReplicationFactor int `mapstructure:"ResponseReplicationFactor"`
	StatusCodeKey     string
	CorrelationIdKey  string
}

type RedisConfig struct {
	Host             string
	Port             string
	User             string
	Pass             string
	ServicesRedisKey string
}

type AuthServer struct {
	AuthenticationKey            string //AUTHENTICATION_KEY
	WebsocketAuthenticationTopic string
	AuthenticationTopic          string
	TTL                          int64
	Timeout                      int64
}

type SentryConfig struct {
	DSN         string
	Environment string
	Debug       bool
}

type Topics struct {
	Health             string
	PublicTopic        string
	AuthTopic          string // AUTHENTICATION_TOPIC
	WebsocketAuthTopic string // WEBSOCKET_AUTHENTICATION_TOPIC
}

type Configuration struct {
	Server     Server
	Redis      RedisConfig
	Kafka      KafkaConfig
	AuthServer AuthServer
	Sentry     SentryConfig
	Topics     Topics
}
