package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// Import our platform package to connect to Redis
	"surge-engine/internal/platform"
)

func main() {
	fmt.Println("👷 Starting the Background Worker...")

	// 1. Connect to the Buffer (same as our other app!)
	rdb, err := platform.NewRedisClient("localhost:6379")
	if err != nil {
		log.Fatalf("🚨 Worker failed to connect to Redis: %v", err)
	}
	fmt.Println("✅ Worker connected to Redis. Waiting for jobs...")

	// 2. Connect to Postgres (The 'Vault')
	// DSN = Data Source Name (The address/username/password)
	dsn := "postgres://user:password@localhost:5432/analytics_db?sslmode=disable"
	db, err := platform.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("🚨 Worker failed to connect to Postgres: %v", err)
	}
	fmt.Println("✅ Worker connected to Postgres!")

	// 3. Create the Table if it doesn't exist (The Architect's Safety Check)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		data JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	ctx := context.Background()

	// 2. The Infinite Loop (The heartbeat of a background worker)
	for {
		// 3. BRPOP: Block and wait on the "analytics" list. 
		// The '0' means "Wait forever, do not timeout".
		result, err := rdb.BRPop(ctx, 0, "analytics").Result()
		
		if err != nil {
			log.Printf("⚠️ Error pulling from queue: %v", err)
			time.Sleep(1 * time.Second) // If Redis hiccups, pause for 1 second so we don't crash the CPU
			continue
		}

		// result is a slice (array) of 2 strings: 
		// result[0] is the name of the list ("analytics")
		// result[1] is the actual JSON payload we pushed earlier
		jsonData := result[1]

		fmt.Printf("📥 BOOM! Worker processed an event: %s\n", jsonData)

		// 3. Save the JSON straight into Postgres!
    	_, err = db.Exec("INSERT INTO events (data) VALUES ($1)", jsonData)
    	if err != nil {
        	log.Printf("❌ Failed to save to DB: %v", err)
    	} else {
        	fmt.Println("💾 Saved to Permanent Vault!")
    	}
	}
}