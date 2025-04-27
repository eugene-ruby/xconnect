package app

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/eugene-ruby/xconnect/rabbitmq"
	"github.com/eugene-ruby/xconnect/rabbitmq/mocks"
)

// TestApplication_PublishMessage verifies that PublishMessage correctly sends a message to the queue.
func TestApplication_PublishMessage(t *testing.T) {
	mock := mocks.NewMockChannel()
	app := NewApplication(mock)

	err := app.PublishMessage("hello world")
	require.NoError(t, err)

	require.Len(t, mock.PublishedMessages, 1)
	require.Equal(t, []byte("hello world"), mock.PublishedMessages[0])
}

// TestApplication_WorkerReceivesMessage verifies that the Worker can consume and process a message.
// Note: This test uses a short sleep and manual context cancellation for shutdown.
// For a more robust version, see TestApplication_WorkerReceivesMessage_WithWaitGroup.
func TestApplication_WorkerReceivesMessage(t *testing.T) {
	mock := mocks.NewMockChannel()

	app := NewApplication(mock)

	// Create a cancellable context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start the worker
	err := app.Start(ctx)
	require.NoError(t, err)

	// Push a test message into the mock consume channel
	mock.ConsumeMessages <- rabbitmq.Delivery{Body: []byte("test message")}

	// Give the worker some time to process the message
	time.Sleep(50 * time.Millisecond)

	// Cancel the context to stop the worker
	cancel()
	app.Wait()
}

// TestApplication_WorkerStartConsumeError verifies that the Worker properly handles Consume errors at startup.
func TestApplication_WorkerStartConsumeError(t *testing.T) {
	mock := mocks.NewMockChannel()
	mock.ConsumeErr = errors.New("failed to consume")

	app := NewApplication(mock)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := app.Start(ctx)
	require.Error(t, err)
}

// TestApplication_WorkerReceivesMessage_WithWaitGroup is a clean and reliable version
// of worker message processing test using sync.WaitGroup without relying on sleeps.
func TestApplication_WorkerReceivesMessage_WithWaitGroup(t *testing.T) {
	mock := mocks.NewMockChannel()

	// Channel to capture received message
	received := make(chan []byte, 1)

	// Use a WaitGroup to wait for the handler to complete
	var wg sync.WaitGroup
	wg.Add(1)

	// Create a custom worker with a custom handler that captures the message and signals completion
	worker := rabbitmq.NewWorker(mock, rabbitmq.WorkerConfig{
		Queue:       "app_queue",
		ConsumerTag: "test_worker",
		AutoAck:     true,
		Handler: func(d rabbitmq.Delivery) error {
			received <- d.Body
			wg.Done()
			return nil
		},
	})

	app := &Application{
		publisher: rabbitmq.NewPublisher(mock),
		worker:    worker,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start the worker
	err := app.Start(ctx)
	require.NoError(t, err)

	// Send a message into the mock channel
	mock.ConsumeMessages <- rabbitmq.Delivery{Body: []byte("test message")}

	// Wait until the message is processed
	wg.Wait()

	// Cancel the context and wait for the worker to fully stop
	cancel()
	app.Wait()

	// Verify that the message was correctly received and processed
	select {
	case msg := <-received:
		require.Equal(t, []byte("test message"), msg)
	default:
		t.Fatal("expected to receive a message but channel was empty")
	}
}
