package rabbitmq_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/eugene-ruby/xconnect/rabbitmq"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

func TestIntegration_WrapAMQPChannel(t *testing.T) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	require.NotEmpty(t, rabbitURL, "RABBITMQ_URL must be set")

	conn, err := amqp.Dial(rabbitURL)
	require.NoError(t, err)
	defer conn.Close()

	rawCh, err := conn.Channel()
	require.NoError(t, err)
	defer rawCh.Close()

	ch := rabbitmq.WrapAMQPChannel(rawCh)

	queueName := "xconnect_integration_test"

	// Declare a queue
	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	require.NoError(t, err)
	require.Equal(t, queueName, q.Name)

	// Publish a test message
	body := []byte("hello integration test")
	err = ch.Publish("", queueName, body)
	require.NoError(t, err)

	// Consume the message
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	require.NoError(t, err)

	select {
	case msg := <-msgs:
		require.Equal(t, body, msg.Body)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestIntegration_Worker(t *testing.T) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	require.NotEmpty(t, rabbitURL, "RABBITMQ_URL must be set")

	conn, err := amqp.Dial(rabbitURL)
	require.NoError(t, err)
	defer conn.Close()

	rawCh, err := conn.Channel()
	require.NoError(t, err)
	defer rawCh.Close()

	channel := rabbitmq.WrapAMQPChannel(rawCh)
	publisher := rabbitmq.NewPublisher(channel)

	queueName := "xconnect_worker_integration_test"

	// Declare a queue
	_, err = channel.QueueDeclare(queueName, false, false, false, false, nil)
	require.NoError(t, err)

	// Channel to capture received messages
	received := make(chan []byte, 1)

	worker := rabbitmq.NewWorker(channel, rabbitmq.WorkerConfig{
		Queue:       queueName,
		ConsumerTag: "worker_test",
		AutoAck:     true,
		Handler: func(d rabbitmq.Delivery) error {
			received <- d.Body
			return nil
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = worker.Start(ctx)
	require.NoError(t, err)

	// Publish a test message
	testBody := []byte("hello worker test")
	err = publisher.Publish("", queueName, testBody)
	require.NoError(t, err)

	select {
	case body := <-received:
		require.Equal(t, testBody, body)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for worker to consume message")
	}

	cancel()
	worker.Wait()
}
