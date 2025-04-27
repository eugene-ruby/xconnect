# Makefile for xconnect project

# Default variables
RABBITMQ_URL ?= amqp://guest:guest@localhost:5672/

.PHONY: all unit-test integration-test run-example docker-up docker-down clean

all: unit-test integration-test

# Run unit tests
unit-test:
	@echo "Running unit tests..."
	go test -v ./rabbitmq/... --timeout 5s

examples-app-test:
	go test -v ./examples/app/... --timeout 10s

# Run integration tests (requires RabbitMQ up)
integration-test:
	@echo "Running integration tests..."
	RABBITMQ_URL=$(RABBITMQ_URL) go test -v ./tests/integration/rabbitmq/... --timeout 5s

docker-up:
	@echo "Starting RabbitMQ with docker-compose..."
	docker-compose -f docker-compose.test.yml up -d

docker-down:
	@echo "Stopping RabbitMQ..."
	docker-compose -f docker-compose.test.yml down

# Run example client
run-example: docker-up
	@echo "Running example client..."
	RABBITMQ_URL=$(RABBITMQ_URL) go run ./examples/rabbitmq/

# Clean temporary files
clean:
	@echo "Cleaning..."
	rm -f coverage.out
