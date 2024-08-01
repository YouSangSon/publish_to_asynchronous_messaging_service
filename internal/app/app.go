package app

import (
	"context"
	"publish_to_messaging_service/internal/app/handler"
	"publish_to_messaging_service/internal/app/kafkaStore"
	"publish_to_messaging_service/internal/config"
	"publish_to_messaging_service/pkg/domain"

	"publish_to_messaging_service/pkg/logger"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"cloud.google.com/go/pubsub"
)

// app is the struct that holds the application
type app struct {
	mode    string
	runMode string
	ctx     context.Context
	topic   interface{}
}

// newApp creates a new instance of the app
func newApp(ctx context.Context, mode, runMode string) (domain.App, error) {
	return &app{
		mode:    mode,
		runMode: runMode,
		ctx:     ctx,
	}, nil
}

// appInit is the first function to be called when the application starts
func (ap *app) Init() error {
	logger.Sugar.Info("app is initializing")

	var newClient interface{}
	var err error

	switch {
	case ap.runMode == "pubsub":
		client, err := pubsub.NewClient(ap.ctx, config.ModeConfig.PubSubs.Project)
		if err != nil {
			return err
		}

		if ok, err := client.Topic(config.ModeConfig.PubSubs.Topic).Exists(ap.ctx); !ok {
			if err != nil {
				return err
			}

			newClient, err = client.CreateTopic(ap.ctx, config.ModeConfig.PubSubs.Topic)
			if err != nil {
				return err
			}
		} else {
			newClient = client.Topic(config.ModeConfig.PubSubs.Topic)
		}
		newClient.(*pubsub.Topic).EnableMessageOrdering = true

		ap.topic = newClient
	case ap.runMode == "kafka":
		newClient, err = kafkaStore.NewKafkaStore(ap.ctx, config.ModeConfig.KafkaStore.Host, config.ModeConfig.KafkaStore.Port)
		if err != nil {
			return err
		}

		ap.topic = newClient
	}

	return nil
}

// Application is the function that runs the application
func Application(ctx context.Context, mode, runMode string) {
	serviceApp, err := newApp(ctx, mode, runMode)
	if err != nil {
		panic(err)
	}

	if err := serviceApp.Init(); err != nil {
		panic(err)
	}

	if err := serviceApp.Run(); err != nil {
		logger.Sugar.Panic(err)
	}
}

// AppRun is the function that runs the application
func (a *app) Run() error {
	fiberApp := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	handlerContext := &handler.HandlerContext{
		RunMode: a.runMode,
		Fiber:   fiberApp,
		Topic:   a.topic,
	}

	fiberApp.Use(fiberLogger.New())
	fiberApp.Use(recover.New())
	fiberApp.Use(limiter.New(limiter.Config{
		Max:        10000,
		Expiration: 30 * time.Second,
	}))
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://example.com, https://example.com",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }))
	// app.Use(compress.New(compress.Config{
	// 	Level: compress.LevelBestSpeed,
	// }))
	// app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: []byte("secret"),
	// }))

	v1Server(handlerContext)

	return fiberApp.Listen(":3000")
}

// v1Server is the function that handles the v1 server
func v1Server(handlerContext *handler.HandlerContext) {
	v1 := handlerContext.Fiber.Group("/v1")

	data := v1.Group("/data-type")
	{
		data.Post("/data", handlerContext.PostData)
	}
}
