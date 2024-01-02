package cmd

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"repo.abanicon.com/abantheter-microservices/websocket/app/configs"
	"repo.abanicon.com/abantheter-microservices/websocket/app/handler"
	"repo.abanicon.com/abantheter-microservices/websocket/app/messaging"
	"repo.abanicon.com/abantheter-microservices/websocket/app/wiring"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/contracts"
	"repo.abanicon.com/abantheter-microservices/websocket/domain/hub"
	"repo.abanicon.com/abantheter-microservices/websocket/infra/broker"
	"repo.abanicon.com/abantheter-microservices/websocket/infra/healer"
	"repo.abanicon.com/abantheter-microservices/websocket/infra/rate_limiter"
	"repo.abanicon.com/abantheter-microservices/websocket/infra/storage"
	"repo.abanicon.com/abantheter-microservices/websocket/infra/utils"
	"repo.abanicon.com/public-library/glogger"
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
		//for i := 10; i > 0; i-- {
		//	time.Sleep(time.Second * 1)
		//	fmt.Printf("\033[2K\rShutting down ... : %d", i)
		//}

		// Perform any necessary cleanup before exiting
		fmt.Println("\nService exited successfully.")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(ctx context.Context) {
	//go func() {
	//	fmt.Println(http.ListenAndServe("localhost:6000", nil))
	//}()
	initLogging(ctx)
	initSentry()
	initWiring(ctx)

	// create topics and health check

	newHub := wiring.Wiring.Hub
	newKafka := wiring.Wiring.Kafka
	newAuthService := wiring.Wiring.GetAuthService()
	newMonitoring := wiring.Wiring.GetMonitoringService()

	messagingService := messaging.NewMessaging(newKafka, newMonitoring, newHub, newAuthService)
	messagingService.Run(ctx)

	go runHttpServer(ctx, newHub)
}

func initLogging(ctx context.Context) {
	glogger.SetConfig(glogger.Config{Production: configs.Config.IsProduction()})
	//sub := monitoring.NewSentrySubscriber(monitoring.SentrySubscriberConfig{Method: monitoring.Async})
	//sub.Run(ctx)
	//logger.Subscribe(sub)
}

func runHttpServer(ctx context.Context, hub *hub.Hub) {
	// init
	monitoringService := wiring.Wiring.GetMonitoringService()

	newPrivateHandler := handler.NewPrivateHandler(hub, monitoringService)
	privateHandler := handler.LoggerMiddleware(
		handler.PrivateChannelMiddleware(
			handler.SocketValidationMiddleware(newPrivateHandler)))

	newPublicHandler := handler.NewPublicHandler(hub, monitoringService)
	publicHandler := handler.LoggerMiddleware(
		handler.SocketValidationMiddleware(
			handler.RateLimiterMiddleware(newPublicHandler)))

	// build endpoints
	mux := http.NewServeMux()
	mux.Handle("/", publicHandler)
	mux.Handle("/private", privateHandler)
	mux.Handle("/metrics", promhttp.Handler())

	address := net.JoinHostPort(configs.Config.Server.Host, configs.Config.Server.Port)
	logrus.WithField("HTTP_Host", configs.Config.Server.Host).
		WithField("HTTP_Port", configs.Config.Server.Port).
		Info("starting HTTP/REST websocket...")

	server := &http.Server{Addr: address, Handler: mux}

	//go func() {
	//	<-ctx.Done()
	//	server.Shutdown(ctx)
	//}()

	err := server.ListenAndServe()
	if err != nil {
		glogger.Fatal(err)
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

// todo move this ito infra section
func initWiring(ctx context.Context) {
	redisProvider, err := storage.NewRedis(configs.Config.Redis)
	if err != nil {
		glogger.Fatalf("Fatal error on create redis("+configs.Config.Redis.Host+":"+configs.Config.Redis.Port+")connection: %s \n", err)
	}

	kafkaProvider, err := broker.NewKafka(configs.Config.Kafka)
	if err != nil {
		glogger.Fatalf("Fatal error on create kafka connection: %s \n", err)
	}

	rateLimiter := rate_limiter.NewRateLimiter(configs.Config.RateLimiter)

	wiring.Wiring = wiring.NewWire(kafkaProvider, redisProvider, rateLimiter, configs.Config)

	// init healer for services
	infraHealer := healer.NewHealerService([]contracts.IProvider{redisProvider}, 1)
	infraHealer.Start(ctx)

	// register profiling
	utils.NewProfiling(configs.Config.Debug).Register()

	glogger.Info("wiring initialized")
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
		glogger.Fatalf("sentry.Init: %s", err)
	}
}
