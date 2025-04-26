package rabbitmq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeliveryConversion(t *testing.T) {
	rawBody := []byte("test-body")
	routingKey := "test.routing.key"

	delivery := Delivery{
		Body:       rawBody,
		RoutingKey: routingKey,
	}

	require.Equal(t, rawBody, delivery.Body)
	require.Equal(t, routingKey, delivery.RoutingKey)
}

func TestTableConversion(t *testing.T) {
	rawTable := Table{
		"x-message-ttl": int32(60000),
	}

	converted := tableToAMQP(rawTable)

	require.Equal(t, int32(60000), converted["x-message-ttl"])
}
