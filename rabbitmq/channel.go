// Package rabbitmq provides abstractions and wrappers for messaging with RabbitMQ.
package rabbitmq

// Channel defines an abstract message channel interface.

// Channel abstracts a message broker channel.
type Channel interface {
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table) error
	Publish(exchange, routingKey string, body []byte) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args Table) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan Delivery, error)
	Close() error
}

// Queue defines a declared queue.
type Queue struct {
	Name      string
	Messages  int
	Consumers int
}
