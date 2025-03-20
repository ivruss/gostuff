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

type Publisher struct {
	Publisher *rabbitmq.Publisher
}

func NewRabbitPublisher(conn *rabbitmq.Conn, exchangeName string) (*Publisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating rabbitMQ publisher: %w", err)
	}

	return &Publisher{Publisher: publisher}, nil
}

func (p *Publisher) Close() error {
	p.Publisher.Close()
	return nil
}

type Consumer struct {
	Consumer *rabbitmq.Consumer
}

func NewRabbitConsumer(
	conn *rabbitmq.Conn,
	queueName string,
	routingKey string,
	exchangeName string,
) (*Consumer, error) {
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

	return &Consumer{Consumer: consumer}, nil
}

func (cs *Consumer) Close() error {
	cs.Consumer.Close()
	return nil
}
