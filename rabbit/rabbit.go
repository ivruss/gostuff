package rabbit

import (
	"fmt"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

func NewRabbitConnection(connstr string) (*rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(
		connstr,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, fmt.Errorf("error opening rabbitMQ connection: %w", err)
	}

	return conn, nil
}

func NewRabbitPublisher(conn *rabbitmq.Conn, exchangeName string) (*rabbitmq.Publisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating rabbitMQ publisher: %w", err)
	}

	return publisher, nil
}

func NewRabbitConsumer(
	conn *rabbitmq.Conn,
	queueName string,
	routingKey string,
	exchangeName string,
) (*rabbitmq.Consumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		conn,
		queueName,
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsExchangeName(exchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating rabbitMQ consumer: %w", err)
	}

	return consumer, nil
}
