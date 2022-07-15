package client

import (
	"context"
	"errors"
	"fmt"
	"track-my-package/app/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueName = "package_status"
)

type rabbitMqClient struct {
	conn             *amqp.Connection
	ch               *amqp.Channel
	connectionString string
	packageStatus    <-chan amqp.Delivery
}

func NewRabbitMqClient(connectionString string) (*rabbitMqClient, error) {
	c := &rabbitMqClient{
		connectionString: connectionString,
	}
	var err error
	c.conn, err = amqp.Dial(c.connectionString)
	if err != nil {
		return nil, err
	}
	c.ch, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}
	err = c.configureQueue()
	return c, err
}

func (c *rabbitMqClient) ConsumeByVehicleID(ctx context.Context, vehicleID string) ([]byte, error) {
	for msg := range c.packageStatus {
		if msg.MessageId == vehicleID {
			_ = msg.Ack(false)
			return msg.Body, nil
		}
	}
	return nil, errors.New("error: when getting package status on channel")
}
func (c *rabbitMqClient) Publish(p *domain.Package) error {
	jsonStr := fmt.Sprintf(`{ "from": %q, "to": %q, "vehicle_id": %q }`, p.From, p.To, p.VehicleID)
	return c.ch.Publish(
		"",
		QueueName,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			MessageId:   p.VehicleID,
			Body:        []byte(jsonStr),
		})
}

func (c *rabbitMqClient) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *rabbitMqClient) configureQueue() error {
	_, err := c.ch.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}
	c.packageStatus, err = c.ch.Consume(
		QueueName,
		"",
		false,
		false,
		false,
		true,
		nil,
	)
	return err
}
