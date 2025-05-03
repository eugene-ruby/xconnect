package mocks_test

import (
	"testing"

	"github.com/eugene-ruby/xconnect/rabbitmq"
	"github.com/eugene-ruby/xconnect/rabbitmq/mocks"
	"github.com/stretchr/testify/require"
)

func TestMockChannel_PublishCapture(t *testing.T) {
	mock := mocks.NewMockChannel()

	err := mock.Publish("test-exchange", "test.key", []byte("hello world"))
	require.NoError(t, err)

	require.Len(t, mock.PublishedMessages, 1)

	msg := mock.PublishedMessages[0]
	require.Equal(t, "test-exchange", msg.Exchange)
	require.Equal(t, "test.key", msg.RoutingKey)
	require.Equal(t, []byte("hello world"), msg.Body)
}

func TestMockChannel_ConsumeSuccess(t *testing.T) {
	mock := mocks.NewMockChannel()

	// Simulate a pushed message
	expected := rabbitmq.Delivery{Body: []byte("incoming message")}
	mock.ConsumeMessages <- expected

	ch, err := mock.Consume("queue", "consumer", true, false, false, false, nil)
	require.NoError(t, err)

	// Read the message from the returned channel
	msg := <-ch
	require.Equal(t, expected.Body, msg.Body)
}

func TestMockChannel_ConsumeError(t *testing.T) {
	mock := mocks.NewMockChannel()
	mock.ConsumeErr = errFake("consume error")

	ch, err := mock.Consume("queue", "consumer", true, false, false, false, nil)
	require.Nil(t, ch)
	require.EqualError(t, err, "consume error")
}

func TestMockChannel_ExchangeDeclare(t *testing.T) {
	mock := mocks.NewMockChannel()

	err := mock.ExchangeDeclare("ex", "direct", true, false, false, false, nil)
	require.NoError(t, err)
}

func TestMockChannel_QueueDeclare(t *testing.T) {
	mock := mocks.NewMockChannel()

	q, err := mock.QueueDeclare("queue", true, false, false, false, nil)
	require.NoError(t, err)
	require.Equal(t, "queue", q.Name)
}

func TestMockChannel_QueueBind(t *testing.T) {
	mock := mocks.NewMockChannel()

	err := mock.QueueBind("queue", "key", "exchange", false, nil)
	require.NoError(t, err)
}

func TestMockChannel_Close(t *testing.T) {
	mock := mocks.NewMockChannel()

	err := mock.Close()
	require.NoError(t, err)
}

// Helper for fake errors
type errFake string

func (e errFake) Error() string {
	return string(e)
}
