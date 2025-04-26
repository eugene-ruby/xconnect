package rabbitmq

// Worker consumes messages from a queue and processes them with a HandlerFunc.

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// HandlerFunc defines a function to process incoming messages.
type HandlerFunc func(Delivery) error

// WorkerConfig holds configuration for a Worker.
type WorkerConfig struct {
	Queue       string
	ConsumerTag string
	AutoAck     bool
	Handler     HandlerFunc
}

// Worker represents a consumer of messages from a queue.
type Worker struct {
	config WorkerConfig
	channel Channel
	wg      sync.WaitGroup
}

// NewWorker creates a new Worker with the given Channel and configuration.
func NewWorker(channel Channel, config WorkerConfig) *Worker {
	return &Worker{
		config: config,
		channel: channel,
	}
}

// Start begins consuming messages from the configured queue.
func (w *Worker) Start(ctx context.Context) error {
	if w.config.Handler == nil {
		return errors.New("worker: HandlerFunc must be set")
	}

	msgs, err := w.channel.Consume(
		w.config.Queue,
		w.config.ConsumerTag,
		w.config.AutoAck,
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("worker: failed to consume: %w", err)
	}

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				if err := w.config.Handler(msg); err != nil {
					fmt.Printf("worker: handler error: %v\n", err)
				}
			}
		}
	}()

	return nil
}

// Wait blocks until the worker has stopped processing.
func (w *Worker) Wait() {
	w.wg.Wait()
}
