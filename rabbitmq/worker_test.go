package rabbitmq

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorker_StartAndHandleMessages(t *testing.T) {
	messages := make(chan Delivery, 1)
	messages <- Delivery{Body: []byte("test message")}
	close(messages)

	mock := &mockChannel{messages: messages}

	handled := false
	handler := func(d Delivery) error {
		handled = true
		return nil
	}

	worker := NewWorker(mock, WorkerConfig{
		Queue:       "test_queue",
		ConsumerTag: "test_consumer",
		AutoAck:     true,
		Handler:     handler,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := worker.Start(ctx)
	require.NoError(t, err)

	worker.Wait()

	require.True(t, handled)
}

func TestWorker_StartConsumeError(t *testing.T) {
	mock := &mockChannel{consumeErr: errors.New("consume failed")}

	worker := NewWorker(mock, WorkerConfig{
		Queue:       "test_queue",
		ConsumerTag: "test_consumer",
		AutoAck:     true,
		Handler:     func(d Delivery) error { return nil },
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := worker.Start(ctx)
	require.Error(t, err)
}

func TestWorker_ContextCancellation(t *testing.T) {
	messages := make(chan Delivery)
	mock := &mockChannel{messages: messages}

	handled := make(chan struct{})

	handler := func(d Delivery) error {
		handled <- struct{}{}
		return nil
	}

	worker := NewWorker(mock, WorkerConfig{
		Queue:       "test_queue",
		ConsumerTag: "test_consumer",
		AutoAck:     true,
		Handler:     handler,
	})

	ctx, cancel := context.WithCancel(context.Background())

	err := worker.Start(ctx)
	require.NoError(t, err)

	cancel()

	// No messages, but worker should shut down cleanly
	worker.Wait()
}
