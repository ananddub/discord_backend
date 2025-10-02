package main

import (
	"log"

	"discord/config"
	"discord/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
