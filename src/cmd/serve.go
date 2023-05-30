package cmd

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"repo.abanicon.com/abantheter-microservices/websocket/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/broker"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/handler"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/infra/storage"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/rate_limiter"
	"repo.abanicon.com/abantheter-microservices/websocket/pkg/wiring"
	"runtime"
	"time"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serving service.",
	Long:  `Serving service.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())

		serve(ctx)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
		fmt.Println()
		time.Sleep(time.Second * 10)
		// Perform any necessary cleanup before exiting
		fmt.Println("Shutting down...")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(ctx context.Context) {
	//initStatePrinter()
	initSentry()
	initWiring()
	initHealer(ctx)
	createTopics(ctx)
	kafkaHealthCheck(ctx)

	newHub := wiring.Wiring.Hub
	newHub.Streaming(ctx)

	auth := wiring.Wiring.GetNewAuthService()
	go auth.Consume(ctx)

	go runHttpServer(ctx, newHub)

}

func createTopics(ctx context.Context) {
	kafkaConfig := configs.Config.Kafka

	topics := []string{
		configs.Config.Topics.Health,
		configs.Config.Topics.PublicTopic,
	}
	kafkaAdmin := wiring.Wiring.GetKafkaAdmin()
	if err := kafkaAdmin.CreateTopic(ctx, topics, kafkaConfig.Partitions, kafkaConfig.ReplicationFactor); err != nil {
		panic(err)
	}
}

func kafkaHealthCheck(ctx context.Context) {
	healthTopic := configs.Config.Topics.Health
	testingMessage := broker.ProduceMessage{
		Topic:         healthTopic,
		Key:           "health-check-key",
		Message:       "valid-health-check-value",
		ResponseTopic: healthTopic,
	}

	broker.Produce(ctx, testingMessage)

	// Ok till here

	result, err := broker.ConsumeHealth(ctx, healthTopic)
	if err != nil {
		panic(err)
	}
	if result == testingMessage.Message {
		log.Info("Kafka is listening")
	}
}

func runHttpServer(ctx context.Context, hub *hub.Hub) {
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	// init
	newPrivateHandler := handler.NewPrivateHandler(hub)
	privateHandler := handler.LoggerMiddleware(
		handler.PrivateChannelMiddleware(
			handler.SocketValidationMiddleware(newPrivateHandler)))

	newPublicHandler := handler.NewPublicHandler(hub)
	publicHandler := handler.LoggerMiddleware(
		handler.SocketValidationMiddleware(
			handler.RateLimiterMiddleware(newPublicHandler)))

	// build endpoints
	mux := http.NewServeMux()
	mux.Handle("/", publicHandler)
	mux.Handle("/private", privateHandler)

	address := net.JoinHostPort(configs.Config.Server.Host, configs.Config.Server.Port)
	log.WithField("HTTP_Host", configs.Config.Server.Host).
		WithField("HTTP_Port", configs.Config.Server.Port).
		Info("starting HTTP/REST websocket...")

	server := &http.Server{Addr: address, Handler: mux}

	//go func() {
	//	<-ctx.Done()
	//	server.Shutdown(ctx)
	//}()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		sentry.CaptureException(err)
	}

	// sentry just capture the main goroutine panics
	defer func() {
		rErr := recover()
		if rErr != nil {
			fmt.Println("Got Err: ", rErr)
			sentry.CurrentHub().Recover(rErr)
			sentry.Flush(time.Second * 5)
		}
	}()
}

func initWiring() {
	redisProvider, err := storage.NewRedis(configs.Config.Redis)
	if err != nil {
		log.Fatalf("Fatal error on create redis("+configs.Config.Redis.Host+":"+configs.Config.Redis.Port+")connection: %s \n", err)
	}

	kafkaProvider, err := broker.NewKafka(configs.Config.Kafka)
	if err != nil {
		log.Fatalf("Fatal error on create kafka connection: %s \n", err)
	}

	newHub := hub.NewHub(kafkaProvider)

	rateLimiter := rate_limiter.NewRateLimiter(configs.Config.RateLimiter)

	wiring.Wiring = wiring.NewWire(kafkaProvider, redisProvider, newHub, rateLimiter, configs.Config)
	log.Info("wiring initialized")
}

func initSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              configs.Config.Sentry.DSN,
		Environment:      configs.Config.Sentry.Environment,
		Debug:            configs.Config.Sentry.Debug,
		TracesSampleRate: 1.0,
		AttachStacktrace: true,
		ServerName:       "Websocket",
		Release:          "0.1",
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

func initHealer(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		defer fmt.Println("closing healer ...")
		for {
			select {
			case <-ticker.C:
				// ping constantly health of redis
				pingErr := wiring.Wiring.Redis.Ping(ctx).Err()
				if pingErr == nil {
					continue
				}

				if err := wiring.Wiring.SetNewRedisInstance(); err != nil {
					fmt.Errorf(err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func printStats() {
	// For memory
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func initStatePrinter() {
	tiker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-tiker.C
			printStats()
		}
	}()
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
