package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/ekeric13/cron/pkg/cron"
)

func main() {
	// Create a base context that cancels on SIGINT (Ctrl+C)
	baseCtx, stopSignal := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stopSignal()

	// Create a context that will also cancel after 5 seconds
	ctx, cancel := context.WithDeadline(baseCtx, time.Now().Add(5*time.Second))
	defer cancel()

	var counter int
	job := cron.Schedule("* * * * * *").Execute(func(ctx context.Context) {
		counter++
		fmt.Printf("Counter incremented to %d\n", counter)
	}).WithContext(ctx)

	job.Start()

	// Wait for the context to be canceled (either SIGINT received or timeout)
	<-ctx.Done()

	// Check why the context was canceled
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Timeout reached, stopping job")
	} else if ctx.Err() == context.Canceled {
		fmt.Println("\nSIGINT received, job stopped gracefully")
	}

	job.Stop()
}
