package rabbitmq

// Publisher wraps a Channel and provides a high-level API for publishing messages.
type Publisher struct {
	ch Channel
}

// NewPublisher creates a new Publisher from an existing Channel.
func NewPublisher(ch Channel) *Publisher {
	return &Publisher{ch: ch}
}

// Publish sends a message to the given exchange with the given routing key.
func (p *Publisher) Publish(exchange, routingKey string, body []byte) error {
	return p.ch.Publish(exchange, routingKey, body)
}

// Close closes the underlying channel.
func (p *Publisher) Close() error {
	return p.ch.Close()
}
