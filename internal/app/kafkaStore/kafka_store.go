package kafkaStore

import (
	"context"
	"fmt"
	"publish_to_messaging_service/pkg/domain"
	"publish_to_messaging_service/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type kafkaEventStore struct {
	ctx    context.Context
	server string
}

// NewKafkaStore creates a new instance of the kafkaEventStore
func NewKafkaStore(ctx context.Context, host string, port int) (domain.KafkaStore, error) {
	kafkaConfig := &kafkaEventStore{
		ctx:    ctx,
		server: fmt.Sprintf("%v:%v", host, port),
	}

	return kafkaConfig, nil
}

// Publish publishes a message to the kafka topic
func (kes *kafkaEventStore) Publish(topic string, key, value []byte, attribute map[string]string) error {
	producerConfig := &kafka.ConfigMap{"bootstrap.servers": kes.server}
	kp, err := kafka.NewProducer(producerConfig)
	if err != nil {
		return err
	}
	defer kp.Close()

	go func() {
		for e := range kp.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Sugar.Errorf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					logger.Sugar.Debugf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// var headers []kafka.Header

	// for key, value := range attribute {
	// 	headers = append(headers, kafka.Header{Key: key, Value: []byte(value)})
	// }

	err = kp.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Key: key, Value: value}, nil)
	if err != nil {
		return err
	}

	for kp.Flush(3000) > 0 {
		logger.Sugar.Info("Still waiting to flush outstanding messages")
	}

	return nil
}
