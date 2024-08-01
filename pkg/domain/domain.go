package domain

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Response struct {
	Result  string      `json:"result"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type DataStream struct {
	Error error
	Data  kafka.Message
}

type KafkaStore interface {
	Publish(topic string, key, value []byte, attribute map[string]string) error
}

type App interface {
	Init() error
	Run() error
}
