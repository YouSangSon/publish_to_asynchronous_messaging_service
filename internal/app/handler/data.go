package handler

import (
	"context"
	"fmt"
	"net/http"
	"publish_to_messaging_service/internal/config"
	dataModel "publish_to_messaging_service/internal/protobuffer/datamodel"
	"publish_to_messaging_service/pkg/domain"
	"publish_to_messaging_service/pkg/logger"

	"cloud.google.com/go/pubsub"
	"github.com/gofiber/fiber/v2"
)

// PostData is the function that handles the POST request to the /data endpoint
func (hc *HandlerContext) PostData(f *fiber.Ctx) error {
	var sensorData dataModel.SensorData

	if err := f.BodyParser(&sensorData); err != nil {
		logger.Sugar.Errorf("error parsing request body : %s", err.Error())
		return f.Status(http.StatusBadRequest).JSON(domain.Response{
			Result:  "failed",
			Message: "error parsing request body",
		})
	}

	if err := hc.Publish(f.Context(), hc.RunMode, hc.Topic, sensorData.Data); err != nil {
		logger.Sugar.Errorf("error publishing data : %s", err.Error())
		return f.Status(http.StatusInternalServerError).JSON(domain.Response{
			Result:  "failed",
			Message: "error publishing data",
		})
	}

	return f.Status(http.StatusOK).JSON(domain.Response{
		Result:  "success",
		Message: "",
	})
}

// Publish is the function that publishes the data to the messaging service
func (hc *HandlerContext) Publish(ctx context.Context, runMode string, topic interface{}, data []byte) error {
	switch runMode {
	case "pubsub":
		pubsubTopic, ok := topic.(*pubsub.Topic)
		if !ok {
			return fmt.Errorf("invalid topic type for pubsub")
		}
		_, err := pubsubTopic.Publish(ctx, &pubsub.Message{
			Data: data,
		}).Get(ctx)
		if err != nil {
			return fmt.Errorf("pubsub operation failed : %s", err.Error())
		}
	case "kafka":
		KafkaStore, ok := topic.(domain.KafkaStore)
		if !ok {
			return fmt.Errorf("invalid topic type for kafka")
		}

		key := []byte("test")

		if err := KafkaStore.Publish(config.ModeConfig.KafkaStore.Topic, key, data, nil); err != nil {
			return fmt.Errorf("kafka operation failed : %s", err.Error())
		}
	default:
		return fmt.Errorf("unsupported run mode: %s", runMode)
	}

	return nil
}
