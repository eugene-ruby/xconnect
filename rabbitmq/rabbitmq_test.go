package rabbitmq_test

import (
	"testing"

	"github.com/eugene-ruby/xconnect/rabbitmq"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

// MockChannel implements rabbitmq.Channel for testing.
type MockChannel struct {
	Published bool
	Closed    bool
}

func (m *MockChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return nil
}

func (m *MockChannel) Publish(exchange, routingKey string, body []byte) error {
	m.Published = true
	return nil
}

func (m *MockChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{Name: name}, nil
}

func (m *MockChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return nil
}

func (m *MockChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	ch := make(chan amqp.Delivery)
	close(ch)
	return ch, nil
}

func (m *MockChannel) Close() error {
	m.Closed = true
	return nil
}

func TestPublisher_PublishAndClose(t *testing.T) {
	mock := &MockChannel{}
	pub := rabbitmq.NewPublisher(mock)

	err := pub.Publish("test-exchange", "test-key", []byte("hello"))
	require.NoError(t, err)
	require.True(t, mock.Published)

	err = pub.Close()
	require.NoError(t, err)
	require.True(t, mock.Closed)
}
