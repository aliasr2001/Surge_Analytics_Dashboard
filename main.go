package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"surge-engine/internal/platform" // Ensure this matches your 'go mod init' name
	"time"
)

type Event struct {
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// 1. Connect to the Infrastructure (Running in Docker)
	// Note: 'localhost:6379' works because you mapped ports in docker-compose
	rdb, err := platform.NewRedisClient("localhost:6379")
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("✅ Connected to the Redis Buffer!")

	// 2. Create our Analytics Event
	event := Event{
		UserID:    "user_456",
		EventType: "page_view",
		CreatedAt: time.Now(),
	}

	// 3. Convert the struct to JSON (the universal language of APIs)
	payload, _ := json.Marshal(event)

	// 4. Push to the "Tunnel" (Redis List named 'analytics')
	ctx := context.Background()
	err = rdb.LPush(ctx, "analytics", payload).Err()
	if err != nil {
		log.Printf("Failed to push to buffer: %v", err)
	}

	fmt.Println("🚀 Event pushed to the buffer! The pipeline is flowing.")
}