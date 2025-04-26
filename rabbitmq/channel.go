package rabbitmq

import "github.com/streadway/amqp"

// Channel defines an abstract interface over AMQP operations.
type Channel interface {
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	Publish(exchange, routingKey string, body []byte) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Close() error
}

// Publisher provides a high-level API to publish messages.
type Publisher struct {
	ch Channel
}

func NewPublisher(ch Channel) *Publisher {
	return &Publisher{ch: ch}
}

func (p *Publisher) Publish(exchange, routingKey string, body []byte) error {
	return p.ch.Publish(exchange, routingKey, body)
}

func (p *Publisher) Close() error {
	return p.ch.Close()
}
