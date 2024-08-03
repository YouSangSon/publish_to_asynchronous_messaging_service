package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"publish_to_messaging_service/pkg/domain"
	"publish_to_messaging_service/pkg/logger"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmq struct {
	ctx         context.Context
	conn        *amqp.Connection
	dataChannel *amqp.Channel
}

// InitRabbitMQ initializes the rabbitMQ connection
func NewRabbitMQ(ctx context.Context, server, vhost, user, password string) (domain.RabbitMQ, error) {
	connectString := fmt.Sprintf("amqp://%s:%s@%s%s", user, password, server, vhost)
	conn, err := amqp.DialConfig(connectString, amqp.Config{Heartbeat: time.Second * 0})
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	return &rabbitmq{
		ctx:         ctx,
		conn:        conn,
		dataChannel: ch,
	}, nil
}

func (r *rabbitmq) Close() error {
	r.dataChannel.Close()

	if !r.conn.IsClosed() {
		return r.conn.Close()
	}

	return nil
}

func (r *rabbitmq) Subscribe(topic domain.Topic) (<-chan amqp.Delivery, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}

	var msgs <-chan amqp.Delivery

	switch topic {
	case domain.RESULT:
		if err := ch.ExchangeDeclare(
			domain.RESULT.ToString(), // name
			"fanout",                 // type: broadcast
			false,                    // durable
			false,                    // auto-deleted
			false,                    // internal
			false,                    // no-wait
			nil,                      // arguments
		); err != nil {
			return nil, err
		}

		q, err := ch.QueueDeclare(
			"",    // name
			false, // durable
			true,  // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return nil, err
		}

		err = ch.QueueBind(
			q.Name,                   // queue name
			"",                       // routing key
			domain.RESULT.ToString(), // exchange
			false,
			nil)
		if err != nil {
			return nil, err
		}

		msgs, err = ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported topic: %v", topic)
	}

	return msgs, nil
}

func (r *rabbitmq) Publish(topic domain.Topic, data interface{}) error {
	switch topic {
	case domain.DATA:
		if err := r.dataChannel.ExchangeDeclare(
			domain.DATA.ToString(), // name
			"fanout",               // type: routing
			false,                  // durable
			false,                  // auto-deleted
			false,                  // internal
			false,                  // no-wait
			nil,                    // arguments
		); err != nil {
			return err
		}

		b, err := json.Marshal(data)
		if err != nil {
			return err
		}

		logger.Sugar.Debugf("publish data channel : %v\n", string(b))

		if err := r.dataChannel.PublishWithContext(
			r.ctx,
			domain.DATA.ToString(), // exchange
			"",                     // routing key
			false,                  // mandatory
			false,                  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        b,
			}); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported topic: %v", topic)
	}

	return nil
}
