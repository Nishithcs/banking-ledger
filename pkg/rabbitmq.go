package pkg

import (
	"os"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func (r *RabbitMQ) Connect(url string) error {
	var err error
	r.conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}
	r.ch, err = r.conn.Channel()
	return err
}

func (r *RabbitMQ) Publish(queue string, body []byte) error {
	_, err := r.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return r.ch.Publish("", queue, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
}

func (r *RabbitMQ) Consume(queue string) (<-chan []byte, error) {
	_, err := r.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	deliveries, err := r.ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	out := make(chan []byte)
	go func() {
		for msg := range deliveries {
			out <- msg.Body
		}
	}()
	return out, nil
}

func (r *RabbitMQ) Close() error {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}


func AmqpURL() string {
    user := os.Getenv("RABBITMQ_USER")
    password := os.Getenv("RABBITMQ_PASSWORD")
    host := os.Getenv("RABBITMQ_HOST")
    port := os.Getenv("RABBITMQ_PORT")

    return fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
}


