package app

import (
	"context"
	"fmt"
	"github.com/eugene-ruby/xconnect/rabbitmq"
)

// Application defines the main app structure.
type Application struct {
	publisher *rabbitmq.Publisher
	worker    *rabbitmq.Worker
}

// NewApplication initializes the application with Publisher and Worker.
func NewApplication(channel rabbitmq.Channel) *Application {
	publisher := rabbitmq.NewPublisher(channel)

	worker := rabbitmq.NewWorker(channel, rabbitmq.WorkerConfig{
		Queue:       "app_queue",
		ConsumerTag: "app_worker",
		AutoAck:     true,
		Handler: func(d rabbitmq.Delivery) error {
			fmt.Printf("[Handler] Received: %s\n", string(d.Body))
			return nil
		},
	})

	return &Application{
		publisher: publisher,
		worker:    worker,
	}
}

// Start launches the worker.
func (a *Application) Start(ctx context.Context) error {
	return a.worker.Start(ctx)
}

// PublishMessage sends a message to the queue.
func (a *Application) PublishMessage(body string) error {
	return a.publisher.Publish("", "app_queue", []byte(body))
}

// Wait blocks until the worker stops.
func (a *Application) Wait() {
	a.worker.Wait()
}
