package configs

var Config Configuration

type Server struct {
	Host string
	Port string
}

type KafkaConfig struct {
	Host              string
	Port              string
	Group             string
	AuthsGroup        string
	ResponseTopic     string
	Partitions        int `mapstructure:"ResponsePartitionsCount"`
	ReplicationFactor int `mapstructure:"ResponseReplicationFactor"`
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
	AuthenticationKey            string
	WebsocketAuthenticationTopic string
	AuthenticationTopic          string
	AuthenticationPrivateTopic   string
	LoginKey                     string
	LogoutKey                    string
	TTL                          int64
	Timeout                      int64
}

type SentryConfig struct {
	DSN         string
	Environment string
	Debug       bool
}

type Topics struct {
	Health      string
	PublicTopic string
}

type Security struct {
	ValidOrigins []string
}

type RateLimiter struct {
	Active bool
	Rate   int
	Period uint64
}

type Debug struct {
	Active     bool
	GC         bool
	Allocation bool
}

type Pod struct {
	Name string
}

type Configuration struct {
	Pod         Pod
	Server      Server
	Redis       RedisConfig
	Kafka       KafkaConfig
	AuthServer  AuthServer
	Sentry      SentryConfig
	Topics      Topics
	RateLimiter RateLimiter
	Security    Security
	Debug       Debug
}
