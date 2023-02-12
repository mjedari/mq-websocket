package configs

import "os"

func get_env_or_default(_key string, _default string) string {
	_val := os.Getenv(_key)
	if _val == "" {
		_val = _default
	}
	return _val
}

const ENVIRONMENT = "local"
const DEBUG = true
const WebSocketPublicTopic = "websocket_public_topic"
const WebSocketKafkaGroup = "websocket_group"
const AuthenticationKafkaGroup = "authentication_group"
const ReadBufferSize = 1024
const WriteBufferSize = 1024
const WebSocketPort = "8080"
const WEBSOCKET_RESPONSE_KEY = "websocket_response_topic"
const WEBSOCKET_AUTHENTICATION_TOPIC = "websocket_authentication_topic"
const CORRELATION_ID_KEY = "correlation_id"
const AUTHENTICATION_TOPIC = "auth_public"
const AUTHENTICATION_KEY = "get_device_id_by_headers"
const AuthTimeout = 20

const SENTRY_DSN = "https://5f3be1a360174839ae243a1f1344866d@o4503970402992128.ingest.sentry.io/4503970404302848"
