package mocks

import (
	"github.com/eugene-ruby/xconnect/rabbitmq"
)

// PublishedMessage represents a captured Publish call in the mock.
type PublishedMessage struct {
	Exchange   string
	RoutingKey string
	Body       []byte
}

// MockChannel is a mock implementation of rabbitmq.Channel used for unit tests.
type MockChannel struct {
	PublishedMessages []PublishedMessage
	ConsumeMessages   chan rabbitmq.Delivery
	ConsumeErr        error
	PublishErr        error
}

// NewMockChannel creates a new MockChannel instance.
func NewMockChannel() *MockChannel {
	return &MockChannel{
		PublishedMessages: make([]PublishedMessage, 0),
		ConsumeMessages:   make(chan rabbitmq.Delivery, 10),
	}
}

func (m *MockChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args rabbitmq.Table) error {
	return nil
}

func (m *MockChannel) Publish(exchange, routingKey string, body []byte) error {
	if m.PublishErr != nil {
		return m.PublishErr
	}
	m.PublishedMessages = append(m.PublishedMessages, PublishedMessage{
		Exchange:   exchange,
		RoutingKey: routingKey,
		Body:       body,
	})
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
