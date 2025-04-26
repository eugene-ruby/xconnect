package rabbitmq

import "github.com/streadway/amqp"

// amqpChannelWrapper wraps *amqp.Channel to implement the Channel interface.
type amqpChannelWrapper struct {
	raw *amqp.Channel
}

func WrapAMQPChannel(ch *amqp.Channel) Channel {
	return &amqpChannelWrapper{raw: ch}
}

func (a *amqpChannelWrapper) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return a.raw.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (a *amqpChannelWrapper) Publish(exchange, routingKey string, body []byte) error {
	return a.raw.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/octet-stream",
		Body:        body,
	})
}

func (a *amqpChannelWrapper) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return a.raw.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (a *amqpChannelWrapper) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return a.raw.QueueBind(name, key, exchange, noWait, args)
}

func (a *amqpChannelWrapper) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return a.raw.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (a *amqpChannelWrapper) Close() error {
	return a.raw.Close()
}
