package tester

import (
	"fmt"
	"net/http"
	"sync"
	// "time"
)

func main() {
	// We use a WaitGroup to make sure our "test" doesn't finish 
	// before all the requests are sent.
	var wg sync.WaitGroup

	fmt.Println("🚀 Starting Stress Test: Sending 20 requests rapidly...")

	for i := 1; i <= 20; i++ {
		wg.Add(1)

		// This 'go func' starts a "Goroutine" 
		// (A mini-thread so all 20 happen at once!)
		go func(requestNum int) {
			defer wg.Done()

			// Replace with your actual local API URL
			resp, err := http.Get("http://localhost:8080/ingest") 
			if err != nil {
				fmt.Printf("❌ Request %d: Connection Error\n", requestNum)
				return
			}
			
			fmt.Printf("Request %d: Status %s\n", requestNum, resp.Status)
		}(i)
	}

	wg.Wait()
	fmt.Println("🏁 Stress Test Complete.")
}