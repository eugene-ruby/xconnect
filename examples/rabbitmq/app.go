// Package main provides an example application using xconnect's rabbitmq package.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eugene-ruby/xconnect/rabbitmq"
	"github.com/streadway/amqp"
)

// App encapsulates the RabbitMQ connection, channel, publisher, and worker.
type App struct {
	conn      *amqp.Connection    // Raw AMQP connection
	rawCh     *amqp.Channel       // Raw AMQP channel
	channel   rabbitmq.Channel    // Wrapped Channel interface
	publisher *rabbitmq.Publisher // Publisher to send messages
	worker    *rabbitmq.Worker    // Worker to consume messages
}

// NewApp initializes the RabbitMQ connection, channel, publisher, and worker.
func NewApp(rabbitURL string) (*App, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	channel := rabbitmq.WrapAMQPChannel(ch)
	publisher := rabbitmq.NewPublisher(channel)

	// Declare a queue named "example_queue" for publishing and consuming messages.
	_, err = channel.QueueDeclare("example_queue", false, false, false, false, nil)
	if err != nil {
		conn.Close()
		ch.Close()
		return nil, err
	}

	// Create a worker to consume messages from the queue.
	worker := rabbitmq.NewWorker(channel, rabbitmq.WorkerConfig{
		Queue:       "example_queue",
		ConsumerTag: "example_consumer",
		AutoAck:     true,
		Handler: func(d rabbitmq.Delivery) error {
			log.Printf("[Worker] Received: %s", string(d.Body))
			return nil
		},
	})

	return &App{
		conn:      conn,
		rawCh:     ch,
		channel:   channel,
		publisher: publisher,
		worker:    worker,
	}, nil
}

// StartPublisher starts a goroutine that publishes a message to the queue every second.
func (a *App) StartPublisher(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		count := 0

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				count++
				body := fmt.Sprintf("message #%d", count)
				if err := a.publisher.Publish("", "example_queue", []byte(body)); err != nil {
					log.Printf("[Publisher] Failed to publish: %v", err)
				} else {
					log.Printf("[Publisher] Published: %s", body)
				}
			}
		}
	}()
}

// StartWorker starts the worker to consume messages from the queue.
func (a *App) StartWorker(ctx context.Context) error {
	return a.worker.Start(ctx)
}

// Wait blocks until the worker has finished processing.
func (a *App) Wait() {
	a.worker.Wait()
}

// Shutdown cleanly closes the AMQP channel and connection.
func (a *App) Shutdown() {
	a.rawCh.Close()
	a.conn.Close()
}
