# xconnect

[![Go Report Card](https://goreportcard.com/badge/github.com/eugene-ruby/xconnect)](https://goreportcard.com/report/github.com/eugene-ruby/xconnect)  
[![Build Status](https://github.com/eugene-ruby/xconnect/actions/workflows/ci.yml/badge.svg)](https://github.com/eugene-ruby/xconnect/actions)  
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)


**xconnect** is a lightweight Go library providing unified interfaces and wrappers for external connection services like **RabbitMQ**, **Redis**, and more in the future.

The main goal is to abstract connection management and messaging operations, allowing you to **unit test easily** with mocks and **seamlessly integrate** into your production infrastructure.

---

## ✨ Features

- Clean Go interfaces for message brokers
- Simple wrappers over popular libraries (like streadway/amqp)
- Easy mocking for testing
- Support for both **publishers** and **consumers**
- Ready for use in CI/CD and local development

---

## 📦 Installation

```bash
go get github.com/eugene-ruby/xconnect
```

---

## 🚀 Usage Example

### Producer (Publishing Messages)

```go
package main

import (
	"log"
	"github.com/streadway/amqp"
	"github.com/eugene-ruby/xconnect/rabbitmq"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	channel := rabbitmq.WrapAMQPChannel(ch)
	publisher := rabbitmq.NewPublisher(channel)

	err = publisher.Publish("", "test_queue", []byte("Hello from xconnect!"))
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Println("✅ Message published!")
}
```

### Consumer (Receiving Messages)

```go
package main

import (
	"log"
	"time"
	"github.com/streadway/amqp"
	"github.com/eugene-ruby/xconnect/rabbitmq"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	channel := rabbitmq.WrapAMQPChannel(ch)
	queueName := "test_queue"

	msgs, err := channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("Waiting for messages...")
	for msg := range msgs {
		log.Printf("📩 Received: %s", msg.Body)
		// simulate processing time
		time.Sleep(1 * time.Second)
	}
}
```

---

## 🛠 Local Development

### Requirements
- Go 1.21+
- Docker and Docker Compose (for running integration tests)

### Useful Commands

Run unit tests:
```bash
go test ./rabbitmq/...
```

Run integration tests with a real RabbitMQ:
```bash
docker-compose -f docker-compose.test.yml up -d
RABBITMQ_URL=amqp://guest:guest@localhost:5672/ go test ./tests/integration/...
```

Stop RabbitMQ:
```bash
docker-compose -f docker-compose.test.yml down
```

### Project Structure

```
rabbitmq/        # Core RabbitMQ interfaces and wrappers
redis/           # (planned) Redis interfaces
internal/        # Private helpers (optional)
tests/integration/ # Integration tests working with real services
docker-compose.test.yml  # Docker setup for integration testing
```

---

## 🤝 Contributing

Feel free to open issues or submit pull requests!

- Write clear, tested code
- Follow Go idioms (gofmt, go vet)
- Make PRs small and focused

We appreciate every contribution! ❤️

---

## 📄 License

This project is licensed under the [MIT License](/LICENSE).

---

Enjoy connecting with **xconnect**! 🚀