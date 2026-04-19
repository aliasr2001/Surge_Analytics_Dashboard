package main

import (
	"time"
	"context"
	"fmt"
	"log"
	"database/sql"

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

	batchSize := 10
	var batch []string // This is our "Bus" (the slice to hold JSON data)

	// Create a "Ticker" that sends a signal every 5 seconds
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

	ctx := context.Background()

    fmt.Println("👷 Worker started with Batching (10) and Timer (5s)...")


	// 2. The Infinite Loop (The heartbeat of a background worker)
	for {
		// 3. BRPOP: Block and wait on the "analytics" list. 
		// The '0' means "Wait forever, do not timeout".
		// result, err := rdb.BRPop(ctx, 0, "analytics").Result()
		
		// if err != nil {
		// 	log.Printf("⚠️ Error pulling from queue: %v", err)
		// 	time.Sleep(1 * time.Second) // If Redis hiccups, pause for 1 second so we don't crash the CPU
		// 	continue
		// }

		// // result is a slice (array) of 2 strings: 
		// // result[0] is the name of the list ("analytics")
		// // result[1] is the actual JSON payload we pushed earlier
		// jsonData := result[1]

		// fmt.Printf("📥 BOOM! Worker processed an event: %s\n", jsonData)

		// // 3. Save the JSON straight into Postgres!
    	// _, err = db.Exec("INSERT INTO events (data) VALUES ($1)", jsonData)
    	// if err != nil {
        // 	log.Printf("❌ Failed to save to DB: %v", err)
    	// } else {
        // 	fmt.Println("💾 Saved to Permanent Vault!")
    	// }


		// Batching Logic

		// 1. Grab 1 item from Redis
		// result, err := rdb.BRPop(ctx, 0, "analytics").Result()
		// if err != nil {
		// 	log.Printf("Error: %v", err)
		// 	continue
		// }

		// // 2. Add the item to our memory batch
		// batch = append(batch, result[1])
		// fmt.Printf("📝 Added to batch. Current size: %d/%d\n", len(batch), batchSize)

		// // 3. If the batch is full, FLUSH IT to Postgres
		// if len(batch) >= batchSize {
		// 	err := flushToPostgres(db, batch)
		// 	if err != nil {
		// 		log.Printf("❌ Batch Save Failed: %v", err)
		// 	} else {
		// 		fmt.Println("🚀 BATCH SAVED: 10 events persisted in one transaction!")
		// 	}

		// 	// 4. Clear the batch for the next round
		// 	batch = nil 
		// }

		// timer-based flushing orchestration Logic

		select {
        // SIGNAL 1: A new item arrived in Redis
        case <-time.After(100 * time.Millisecond): // We check Redis frequently
            // Use LPop (non-blocking) or a short timeout RPop here for the select pattern
            result, err := rdb.RPop(ctx, "analytics").Result()
            if err == nil {
                batch = append(batch, result)
                fmt.Printf("📝 Added to batch. Size: %d/%d\n", len(batch), batchSize)
            }

            if len(batch) >= batchSize {
                fmt.Println("🚌 Bus full! Flushing...")
                flushToPostgres(db, batch)
                batch = nil
            }

        // SIGNAL 2: The 5-second timer went off!
        case <-ticker.C:
            if len(batch) > 0 {
                fmt.Println("⏰ Timer hit! Flushing partial batch...")
                flushToPostgres(db, batch)
                batch = nil
            }
        }
	}
}


func flushToPostgres(db *sql.DB, batch []string) error {
	// We build a single query like: INSERT INTO events (data) VALUES ($1), ($2), ($3)...
	query := "INSERT INTO events (data) VALUES "
	values := []any{}

	for i, jsonStr := range batch {
		query += fmt.Sprintf("($%d),", i+1)
		values = append(values, jsonStr)
	}

	// Remove the last comma and add it to the DB
	query = query[0 : len(query)-1]

	_, err := db.Exec(query, values...)
	return err
}