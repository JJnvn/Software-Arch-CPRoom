package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	TOKEN        = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGFkbWluLmNvbSIsInJvbGUiOiJBRE1JTiIsImlzcyI6ImNwcm9vbS1hdXRoIiwic3ViIjoiYjA1ODllMzMtZjRhOS00Y2Y4LThiMGEtMDY0MzFmNTNmODgwIiwiZXhwIjoxNzYzMzc2MjY5LCJpYXQiOjE3NjI3NzE0Njl9.IbR0r6Q2mH7j56RdoLGClZDcwuIq4olEzM8V-iWGOag"
	BOOKINGS_URL = "http://localhost:8000/approvals/pending"
	ROOMS_URL    = "http://localhost:8082/rooms"
)

// Stage defines a load stage
type Stage struct {
	duration time.Duration
	target   int
}

// k6-style stages
var stages = []Stage{
	{duration: 1 * time.Minute, target: 10}, // ramp up
	{duration: 8 * time.Minute, target: 10}, // hold steady
	{duration: 1 * time.Minute, target: 0},  // ramp down
}

// simulateRequest sends a GET request to the given endpoint
func simulateRequest(endpoint string, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{Timeout: 120 * time.Second}

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", "Bearer "+TOKEN)

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("[User %d] Error fetching %s: %v\n", id, endpoint, err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Printf("[User %d] GET %s -> %d\n", id, endpoint, res.StatusCode)
	_ = body

	time.Sleep(5 * time.Second) // simulate user wait between calls
}

// runLoadTest runs the load test for a single endpoint
func runLoadTest(name, endpoint string) {
	fmt.Printf("\n=== Starting load test: %s ===\n", name)
	var wg sync.WaitGroup

	for _, stage := range stages {
		fmt.Printf("\nStage: %v users for %v\n", stage.target, stage.duration)
		for i := 1; i <= stage.target; i++ {
			wg.Add(1)
			go simulateRequest(endpoint, i, &wg)
			time.Sleep(stage.duration / time.Duration(stage.target)) // smooth ramp
		}
		time.Sleep(stage.duration)
	}

	wg.Wait()
	fmt.Printf("=== Load test complete: %s ===\n\n", name)
}

func main() {
	// Run bookings load test
	runLoadTest("Get All Bookings", BOOKINGS_URL)

	// Run rooms load test
	runLoadTest("Get All Rooms", ROOMS_URL)
}
