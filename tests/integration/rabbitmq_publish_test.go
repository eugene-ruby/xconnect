package integration

import (
	"os"
	"testing"
	"time"

	"github.com/eugene-ruby/xconnect/rabbitmq"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

func TestPublisher_Integration(t *testing.T) {
	url := os.Getenv("RABBITMQ_URL")
	require.NotEmpty(t, url, "RABBITMQ_URL environment variable must be set")

	conn, err := amqp.Dial(url)
	require.NoError(t, err)
	defer conn.Close()

	amqpCh, err := conn.Channel()
	require.NoError(t, err)
	defer amqpCh.Close()

	channel := rabbitmq.WrapAMQPChannel(amqpCh)

	pub := rabbitmq.NewPublisher(channel)

	queueName := "test_queue_xconnect"

	_, err = channel.QueueDeclare(queueName, false, false, false, false, nil)
	require.NoError(t, err)

	body := "hello from rabbitmq integration test"

	err = pub.Publish("", queueName, []byte(body))
	require.NoError(t, err)

	msgs, err := channel.Consume(queueName, "", true, false, false, false, nil)
	require.NoError(t, err)

	select {
	case msg := <-msgs:
		require.Equal(t, body, string(msg.Body))
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for message")
	}

	err = pub.Close()
	require.NoError(t, err)
}
