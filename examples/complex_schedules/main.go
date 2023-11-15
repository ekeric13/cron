package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ekeric13/cron/pkg/cron"
)

func main() {
	task := func(ctx context.Context) {
		fmt.Println("Task executed at:", time.Now().Format(time.RFC1123))
	}

	// Job 1: Run once a day at 9 AM PST
	job1 := cron.Schedule("0 9 * * *").Execute(task).SetTimezone(time.FixedZone("PST", -8*3600))

	// Job 2: Run on Tuesdays at noon PST
	job2 := cron.Schedule("0 12 * * 2").Execute(task).SetTimezone(time.FixedZone("PST", -8*3600))

	// Job 3: Run on the 5th day of the month at 8 PM PST
	job3 := cron.Schedule("0 20 5 * *").Execute(task).SetTimezone(time.FixedZone("PST", -8*3600))

	// Start all jobs
	job1.Start()
	job2.Start()
	job3.Start()

	fmt.Println("Jobs started. Press Ctrl+C to exit.")

	// Keep the application running
	select {}
}
