package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"log"
	"time"
)

// HandlerFunc defines a function to process incoming messages.
type HandlerFunc func(Delivery) error

// WorkerConfig holds configuration for a Worker.
type WorkerConfig struct {
	Queue       string
	ConsumerTag string
	AutoAck     bool
	Handler     HandlerFunc

	// automatically declare & bind queue before consuming
	Declare        bool
	BindRoutingKey string
	BindExchange   string
}

// Worker represents a consumer of messages from a queue.
type Worker struct {
	config  WorkerConfig
	channel Channel
	wg      sync.WaitGroup
}

// DeclareAndBind declares and binds the queue with given parameters.
func DeclareAndBind(channel Channel, queue, key, exchange string) error {
	_, err := channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return channel.QueueBind(queue, key, exchange, false, nil)
}

// NewWorker creates a new Worker with the given Channel and configuration.
func NewWorker(channel Channel, config WorkerConfig) *Worker {
	return &Worker{
		config:  config,
		channel: channel,
	}
}

// Start begins consuming messages from the configured queue.
func (w *Worker) Start(ctx context.Context) error {
	if w.config.Handler == nil {
		return errors.New("worker: HandlerFunc must be set")
	}

	// Optional queue declare & bind
	if w.config.Declare {
		if err := DeclareAndBind(w.channel, w.config.Queue, w.config.BindRoutingKey, w.config.BindExchange); err != nil {
			return fmt.Errorf("worker: declare and bind failed: %w", err)
		}
	}

	// Start consuming
	msgs, err := w.channel.Consume(
		w.config.Queue,
		w.config.ConsumerTag,
		w.config.AutoAck,
		false, false, false, nil,
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
				log.Println("👷 canceling consumer", w.config.ConsumerTag)
				_ = w.channel.Cancel(w.config.ConsumerTag, false)
				time.Sleep(200 * time.Millisecond)
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
