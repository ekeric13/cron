package cron

import (
	"context"
	"testing"
	"time"
)

// TestSchedule tests the Schedule function for correct schedule parsing.
func TestSchedule(t *testing.T) {
	// Test with a valid schedule string
	job := Schedule("*/5 * * * * *")
	if job == nil {
		t.Errorf("Schedule returned nil for a valid cron string")
	}

	// Test with an invalid schedule string
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Schedule did not panic with an invalid cron string")
		}
	}()
	Schedule("invalid-cron-string")
}

// TestSetBlocking tests the SetBlocking method.
func TestSetBlocking(t *testing.T) {
	job := Schedule("*/5 * * * * *")

	job.SetBlocking(true)
	if job.Blocking != true {
		t.Errorf("SetBlocking(true) did not set job.Blocking to true")
	}

	job.SetBlocking(false)
	if job.Blocking != false {
		t.Errorf("SetBlocking(false) did not set job.Blocking to false")
	}
}

// TestWithContext tests the WithContext method.
func TestWithContext(t *testing.T) {
	job := Schedule("*/5 * * * * *")
	ctx, cancelFunc := context.WithCancel(context.Background())

	job.WithContext(ctx)

	// Cancel the original context.
	cancelFunc()

	// The derived context should also be canceled.
	select {
	case <-job.Ctx.Done():
		// Test passes: job.Ctx is canceled.
	case <-time.After(1 * time.Millisecond):
		t.Errorf("WithContext did not set the correct context: context is not canceled")
	}
}

// TestSetTimezone tests the SetTimezone method.
func TestSetTimezone(t *testing.T) {
	job := Schedule("*/5 * * * * *")
	loc, _ := time.LoadLocation("America/New_York")

	job.SetTimezone(loc)
	if job.Timezone != loc {
		t.Errorf("SetTimezone did not set the correct timezone")
	}
}

// TestJobExecution tests if a job increments a counter as expected.
func TestJobExecution(t *testing.T) {
	var counter int
	job := Schedule("* * * * * *").Execute(func(ctx context.Context) {
		counter++
	})

	job.Start()

	// Let the job run for slightly more than 1 second.
	time.Sleep(1100 * time.Millisecond)

	job.Stop()

	// We expect the counter to be at least 1 since the job runs every second.
	if counter < 1 {
		t.Errorf("Expected counter to be incremented, got %d", counter)
	}
}
