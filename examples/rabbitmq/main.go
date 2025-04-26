package main

import (
	"context"
	"log"
	"os"
)

func main() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("RABBITMQ_URL environment variable must be set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := NewApp(url)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer app.Shutdown()

	go app.StartPublisher(ctx)

	if err := app.StartWorker(ctx); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	WaitForShutdown(cancel)
	app.Wait()
}
