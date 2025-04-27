package mocks

import (
	"github.com/eugene-ruby/xconnect/rabbitmq"
)

// MockChannel is a mock implementation of rabbitmq.Channel used for unit tests.
type MockChannel struct {
	PublishedMessages [][]byte
	ConsumeMessages   chan rabbitmq.Delivery
	ConsumeErr        error
}

// NewMockChannel creates a new MockChannel instance.
func NewMockChannel() *MockChannel {
	return &MockChannel{
		ConsumeMessages: make(chan rabbitmq.Delivery),
	}
}

func (m *MockChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args rabbitmq.Table) error {
	return nil
}

func (m *MockChannel) Publish(exchange, routingKey string, body []byte) error {
	m.PublishedMessages = append(m.PublishedMessages, body)
	return nil
}

func (m *MockChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args rabbitmq.Table) (rabbitmq.Queue, error) {
	return rabbitmq.Queue{Name: name}, nil
}

func (m *MockChannel) QueueBind(name, key, exchange string, noWait bool, args rabbitmq.Table) error {
	return nil
}

func (m *MockChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args rabbitmq.Table) (<-chan rabbitmq.Delivery, error) {
	if m.ConsumeErr != nil {
		return nil, m.ConsumeErr
	}
	return m.ConsumeMessages, nil
}

func (m *MockChannel) Close() error {
	return nil
}
