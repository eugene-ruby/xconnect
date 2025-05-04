package rabbitmq

// mockChannel is a mock implementation of the Channel interface.

type mockChannel struct {
	messages   chan Delivery
	consumeErr error
	published  bool
	closed     bool
	publishErr error
	closeErr   error
	CancelCalled bool
	CancelArgs   []string
}

func (m *mockChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table) error {
	return nil
}

func (m *mockChannel) Publish(exchange, routingKey string, body []byte) error {
	m.published = true
	return m.publishErr
}

func (m *mockChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (Queue, error) {
	return Queue{Name: name}, nil
}

func (m *mockChannel) QueueBind(name, key, exchange string, noWait bool, args Table) error {
	return nil
}

func (m *mockChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan Delivery, error) {
	if m.consumeErr != nil {
		return nil, m.consumeErr
	}
	return m.messages, nil
}

func (m *mockChannel) Close() error {
	m.closed = true
	return m.closeErr
}

func (m *mockChannel) Cancel(consumer string, noWait bool) error {
	m.CancelCalled = true
	m.CancelArgs = append(m.CancelArgs, consumer)
	return nil
}

var _ Channel = (*mockChannel)(nil)
