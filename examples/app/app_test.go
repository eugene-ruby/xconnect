// File: examples/app/app_test.go
package app

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/require"
// 	"github.com/eugene-ruby/xconnect/rabbitmq"
// )
// // "github.com/eugene-ruby/xconnect/rabbitmq/mocks"

// // TestApplication_PublishMessage verifies that PublishMessage correctly sends a message.
// func TestApplication_PublishMessage(t *testing.T) {
// 	mock := mocks.NewMockChannel()
// 	app := NewApplication(mock)

// 	err := app.PublishMessage("hello world")
// 	require.NoError(t, err)
// 	require.Len(t, mock.publishedMessages, 1)
// 	require.Equal(t, []byte("hello world"), mock.publishedMessages[0])
// }

// // TestApplication_WorkerReceivesMessage verifies that the Worker can consume and handle a message.
// func TestApplication_WorkerReceivesMessage(t *testing.T) {
// 	mock := mocks.NewMockChannel()

// 	// Push a test message into the mocked consume channel
// 	mock.ConsumeMessages <- rabbitmq.Delivery{Body: []byte("test message")}
// 	close(mock.ConsumeMessages)

// 	app := NewApplication(mock)

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	err := app.Start(ctx)
// 	require.NoError(t, err)

// 	app.Wait()
// }

// // TestApplication_WorkerStartConsumeError checks that the application handles Consume errors correctly.
// func TestApplication_WorkerStartConsumeError(t *testing.T) {
// 	mock := mocks.NewMockChannel()
// 	mock.ConsumeErr = errors.New("failed to consume")

// 	app := NewApplication(mock)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	err := app.Start(ctx)
// 	require.Error(t, err)
// }
