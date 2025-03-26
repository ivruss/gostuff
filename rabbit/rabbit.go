package rabbit

import (
	"fmt"
	rabbitmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
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

type Publisher interface {
	Publish(data []byte) error
	Close() error
}

type RabbitPublisher struct {
	publisher    *rabbitmq.Publisher
	exchangeName string
	routingKey   []string
	logger       *zap.Logger
}

func NewRabbitPublisher(conn *rabbitmq.Conn, exchangeName string, routingKey []string, logger *zap.Logger) (Publisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating rabbitMQ publisher: %w", err)
	}

	return &RabbitPublisher{
		publisher:    publisher,
		exchangeName: exchangeName,
		routingKey:   routingKey,
		logger:       logger,
	}, nil
}

func (p *RabbitPublisher) Publish(data []byte) error {
	err := p.publisher.Publish(
		data,
		p.routingKey,
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(p.exchangeName),
	)
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}

	return nil
}

func (p *RabbitPublisher) Close() error {
	p.publisher.Close()
	return nil
}

type Consumer interface {
	Consume(handler func([]byte) error) error
	Close() error
}

type RabbitConsumer struct {
	consumer *rabbitmq.Consumer
	logger   *zap.Logger
}

func NewRabbitConsumer(
	conn *rabbitmq.Conn,
	queueName string,
	routingKey string,
	exchangeName string,
	logger *zap.Logger,
) (Consumer, error) {
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

	return &RabbitConsumer{
		consumer: consumer,
		logger:   logger,
	}, nil
}

func (c *RabbitConsumer) Consume(handler func([]byte) error) error {
	err := c.consumer.Run(
		func(delivery rabbitmq.Delivery) rabbitmq.Action {
			c.logger.Sugar().Infof("recieved rabbit message: %s", string(delivery.Body))

			if err := handler(delivery.Body); err != nil {
				c.logger.Sugar().Errorf("error handling rabbit message: %s", err)
				return rabbitmq.NackRequeue
			}

			return rabbitmq.Ack
		})
	if err != nil {
		return fmt.Errorf("error consuming message: %w", err)
	}
	return nil
}

func (c *RabbitConsumer) Close() error {
	c.consumer.Close()
	return nil
}
