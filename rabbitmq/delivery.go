package rabbitmq

// Delivery represents a message delivered from the broker.
type Delivery struct {
	Body       []byte
	RoutingKey string
}
