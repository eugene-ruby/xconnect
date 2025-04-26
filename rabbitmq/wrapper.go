package rabbitmq

import (
	"github.com/streadway/amqp"
)

// amqpChannelWrapper wraps *amqp.Channel to implement Channel.
type amqpChannelWrapper struct {
	raw *amqp.Channel
}

// WrapAMQPChannel wraps a raw amqp.Channel into Channel.
func WrapAMQPChannel(ch *amqp.Channel) Channel {
	return &amqpChannelWrapper{raw: ch}
}

func (a *amqpChannelWrapper) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table) error {
	return a.raw.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, tableToAMQP(args))
}

func (a *amqpChannelWrapper) Publish(exchange, routingKey string, body []byte) error {
	return a.raw.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/octet-stream",
		Body:        body,
	})
}

func (a *amqpChannelWrapper) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (Queue, error) {
	q, err := a.raw.QueueDeclare(name, durable, autoDelete, exclusive, noWait, tableToAMQP(args))
	if err != nil {
		return Queue{}, err
	}
	return Queue{Name: q.Name, Messages: q.Messages, Consumers: q.Consumers}, nil
}

func (a *amqpChannelWrapper) QueueBind(name, key, exchange string, noWait bool, args Table) error {
	return a.raw.QueueBind(name, key, exchange, noWait, tableToAMQP(args))
}

func (a *amqpChannelWrapper) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan Delivery, error) {
	rawChan, err := a.raw.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, tableToAMQP(args))
	if err != nil {
		return nil, err
	}

	wrappedChan := make(chan Delivery)
	go func() {
		defer close(wrappedChan)
		for msg := range rawChan {
			wrappedChan <- Delivery{
				Body:       msg.Body,
				RoutingKey: msg.RoutingKey,
			}
		}
	}()
	return wrappedChan, nil
}

func (a *amqpChannelWrapper) Close() error {
	return a.raw.Close()
}

func tableToAMQP(t Table) amqp.Table {
	if t == nil {
		return nil
	}
	return amqp.Table(t)
}
