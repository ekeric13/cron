package main

import (
	"context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/ekeric13/cron/pkg/cron"
)

func main() {
	var counter int
	job := cron.Schedule("* * * * * *").Execute(func(ctx context.Context) {
		counter++
		fmt.Printf("Counter incremented to %d\n", counter)
	})

	// Marshal the job to JSON
	jobJSON, err := json.Marshal(job)
	if err != nil {
		fmt.Println("Error marshalling job:", err)
		return
	}

	// Print the JSON representation of the job
	fmt.Println(string(jobJSON))

	job.Start()

	// Let the job run for 5 seconds
	time.Sleep(5 * time.Second)

	job.Stop()

	fmt.Println("Job stopped")
}
