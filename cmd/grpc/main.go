package main

import (
	"discord/internal/app"
	"log"
)

// main is the entry point for the Discord gRPC server
func main() {
	app := app.Application{}

	// Initialize application
	if err := app.Initialize(); err != nil {
		log.Fatalf("❌ Failed to initialize application: %v", err)
	}
	defer app.Shutdown()

	// Start gRPC server
	if err := app.StartServer(); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
