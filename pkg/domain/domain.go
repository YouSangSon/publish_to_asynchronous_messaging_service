package domain

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	amqp "github.com/rabbitmq/amqp091-go"
)

// *---------------------------------------------------------------------------------------------------------------//
type Topic int

const (
	DATA Topic = iota
	RESULT
)

func (t Topic) ToString() string {
	switch t {
	case RESULT:
		return "RESULT"
	case DATA:
		return "DATA"
	default:
		return "UNKNOWN"
	}
}

type RabbitMQ interface {
	Close() error
	Subscribe(topic Topic) (<-chan amqp.Delivery, error)
	Publish(topic Topic, data interface{}) error
}

// *---------------------------------------------------------------------------------------------------------------//
type Response struct {
	Result  string      `json:"result"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// *---------------------------------------------------------------------------------------------------------------//
type DataStream struct {
	Error error
	Data  kafka.Message
}

// *---------------------------------------------------------------------------------------------------------------//
type KafkaStore interface {
	Publish(topic string, key, value []byte, attribute map[string]string) error
}

// *---------------------------------------------------------------------------------------------------------------//
type App interface {
	Init() error
	Run() error
}

// *---------------------------------------------------------------------------------------------------------------//
