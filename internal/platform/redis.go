package platform

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient returns a pointer to a Redis client, and an error if it fails.
// Look at the return type: (*redis.Client, error) -> Returning TWO things!
func NewRedisClient(addr string) (*redis.Client, error) {
	
	// 1. We use our shortcut := to create the client variable
	client := redis.NewClient(&redis.Options{
		Addr:     addr, 
		Password: "", // No password for our local Docker setup
		DB:       0,  // Default database
	})

	// 2. We use := again to get the status of our "Ping" to the database
	status := client.Ping(context.Background())

	// 3. The classic Go error check! Is the error NOT nil?
	if status.Err() != nil {
		return nil, status.Err() // Return "nothing" for the client, and the actual error
	}

	// 4. Success! Return the client, and "nil" (no error)
	return client, nil
}