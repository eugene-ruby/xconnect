# About xConnect
```bash
╭──────────────────────────────────────────╮
│    xconnect: clean abstraction layer      │
│    for RabbitMQ, Redis and more           │
╰──────────────────────────────────────────╯
      ↳ lightweight, testable, production-ready
```

[![Go Report Card](https://goreportcard.com/badge/github.com/eugene-ruby/xconnect)](https://goreportcard.com/report/github.com/eugene-ruby/xconnect)
[![Build Status](https://github.com/eugene-ruby/xconnect/actions/workflows/ci.yml/badge.svg)](https://github.com/eugene-ruby/xconnect/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)


**xconnect** is a lightweight Go library providing unified interfaces and wrappers for external connection services like **RabbitMQ**, **Redis** (planned), and more.

The main goal is to abstract connection management and messaging operations, allowing you to **unit test easily** with mocks and **seamlessly integrate** into production systems.

---

## 🚀 Why Use xconnect/rabbitmq

- Clean separation between transport and business logic.
- Easy-to-mock interfaces for unit tests without real brokers.
- Flexible architecture supporting both low-level access (publish/consume) and high-level Worker abstraction.
- Context-driven cancellation and graceful shutdowns.
- No external types leak into your application domain.
- Designed for real production systems with CI/CD in mind.

---

## ✨ Features

- Clean Go interfaces for messaging and data brokers
- Simple wrappers over popular libraries (e.g., streadway/amqp)
- Easy mocking for unit testing
- Support for both **publishers** and **consumers**
- High-level `Worker` abstraction for consuming queues elegantly
- Ready for CI/CD integration and local development
- No external types leaking through interfaces

---

## 📦 Installation

```bash
go get github.com/eugene-ruby/xconnect
```

---

## 🚀 Usage Examples

See full working examples in [`examples/rabbitmq`](./examples/rabbitmq).

### Producer (Publishing Messages)

```go
// ... your producer code (unchanged) ...
```

### Consumer (Receiving Messages)

```go
// ... your consumer code (unchanged) ...
```

---

## 🧪 Mock Support for Unit Testing

`xconnect` provides ready-to-use mocks for unit testing your applications without requiring a live RabbitMQ server.

- Available under [`rabbitmq/mocks`](./rabbitmq/mocks/).
- Provides `MockChannel` implementing `rabbitmq.Channel` interface.
- Useful for testing publishers, workers, and custom logic.

Example:

```go
import "github.com/eugene-ruby/xconnect/rabbitmq/mocks"

func TestPublish(t *testing.T) {
    mock := mocks.NewMockChannel()
    publisher := rabbitmq.NewPublisher(mock)

    err := publisher.Publish(\"\", \"queue\", []byte(\"test message\"))
    require.NoError(t, err)
}
```

✅ No need for real broker  
✅ Instant feedback during development

---

## 📚 Full Example Applications

- [`examples/rabbitmq`](./examples/rabbitmq) — Basic Producer + Worker example with graceful shutdown.
- [`examples/app`](./examples/app) — Full example of building an application using Publisher, Worker, and unit tests based on mocks.

```bash
RABBITMQ_URL=amqp://guest:guest@localhost:5672/ go run ./examples/rabbitmq/
2025/04/26 16:52:39 [Publisher] Published: message #1
2025/04/26 16:52:39 [Worker] Received: message #1
2025/04/26 16:52:40 [Publisher] Published: message #2
2025/04/26 16:52:40 [Worker] Received: message #2
2025/04/26 16:52:41 [Publisher] Published: message #3
2025/04/26 16:52:41 [Worker] Received: message #3
.....
```

---

## 🛠 Local Development

### Requirements
- Go 1.21+
- Docker and Docker Compose (for integration testing)

### Useful Commands

Run unit tests:

```bash
make unit-test
```

Run integration tests (requires running RabbitMQ):

```bash
make integration-test
```

Run example client:

```bash
make run-example
```

Start RabbitMQ with Docker Compose:

```bash
make docker-up
```

Stop RabbitMQ:

```bash
make docker-down
```

---

## 📂 Project Structure

```
/rabbitmq/               # Core RabbitMQ interfaces and wrappers
/rabbitmq/mocks/         # Public mocks for unit testing (MockChannel, etc.)
/examples/rabbitmq/      # Basic live example (Publisher + Worker + graceful shutdown)
/examples/app/           # Full application example with unit tests using mocks
/tests/integration/      # Integration tests for RabbitMQ (real broker tests)
/cmd/                    # (Reserved for CLI applications if needed)
/internal/               # (Reserved for internal utilities)
/docker-compose.test.yml # Docker Compose setup for integration testing
/go.mod                  # Go module definition
/go.sum                  # Go module checksum file
/README.md               # Project documentation
```

---

## 📜 Worker Overview

### General structure of `Worker` in `xconnect/rabbitmq`

#### 1. Interfaces and types

- `Channel`: abstraction for working with queues and messages.
- `Delivery`: structure describing a received message.
- `HandlerFunc func(Delivery) error`: function for handling messages.
- `Worker`: structure that manages subscription and message processing.

#### 2. What does `Worker` do?

- Subscribes to a queue using `Channel.Consume`.
- Starts a goroutine to read and handle messages via `HandlerFunc`.
- Listens for cancellation via `context.Context`.
- Waits for graceful shutdown using `sync.WaitGroup`.

#### 3. Worker lifecycle

```
NewWorker() -> Start(ctx) -> Consume(queue) -> for msg in msgs -> HandlerFunc(msg) -> Wait()
```

---

## 🤝 Contributing

We welcome contributions!  
Please open issues or submit pull requests.

- Follow Go best practices (`gofmt`, `go vet`)
- Write clear, well-tested code
- Keep pull requests small and focused

---

## 📄 License

This project is licensed under the [MIT License](/LICENSE).
