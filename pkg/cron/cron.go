// Package cron provides a simple yet powerful cron scheduling library.
// It allows users to schedule functions to be executed at specific times or intervals,
// following a format similar to traditional UNIX cron, but with extended support for seconds and milliseconds.
// Parsing of UNIX cron is done by https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.1
package cron

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	_cron "github.com/robfig/cron/v3"
)

// Job represents a cron job with a specific schedule and task.
// It holds the schedule string, the parsed schedule, execution settings like blocking behavior, timezone,
// and the function to execute.
type Job struct {
	scheduleStr string
	Schedule    _cron.Schedule  `json:"schedule"`
	Blocking    bool            `json:"blocking"`
	Timezone    *time.Location  `json:"timezone"`
	Ctx         context.Context `json:"-"`
	cancelFunc  context.CancelFunc
	Fn          func(ctx context.Context) `json:"-"`
	isRunning   bool
	mutex       sync.RWMutex
}

// MarshalJSON customizes the JSON output of Job.
func (j *Job) MarshalJSON() ([]byte, error) {
	type Alias Job
	return json.Marshal(&struct {
		ScheduleStr string `json:"schedule_str"`
		*Alias
	}{
		ScheduleStr: j.scheduleStr,
		Alias:       (*Alias)(j),
	})
}

// Schedule initializes a new Job with a given cron schedule string.
// The function panics if the schedule string is invalid.
// The schedule string supports the traditional UNIX cron format with optional seconds field at the beginning.
func Schedule(scheduleStr string) *Job {
	fields := strings.Fields(scheduleStr)
	var parser _cron.Parser

	if len(fields) == 6 {
		parser = _cron.NewParser(_cron.Second | _cron.Minute | _cron.Hour | _cron.Dom | _cron.Month | _cron.Dow)
	} else {
		parser = _cron.NewParser(_cron.Minute | _cron.Hour | _cron.Dom | _cron.Month | _cron.Dow)
	}

	schedule, err := parser.Parse(scheduleStr)
	if err != nil {
		panic("invalid cron schedule")
	}
	// Default context
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Job{
		scheduleStr: scheduleStr,
		Schedule:    schedule,
		// Default non-blocking
		Blocking: false,
		// Default to UTC
		Timezone:   time.UTC,
		Ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}

// now returns the current time in the Job's timezone.
func (j *Job) now() time.Time {
	return time.Now().In(j.Timezone)
}

// SetBlocking configures the Job's blocking behavior.
// If set to true, the job will run its task synchronously. If false, the job will run asynchronously.
func (j *Job) SetBlocking(blocking bool) *Job {
	// locking in case you change on the fly but would not recommend
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.Blocking = blocking
	return j
}

// WithContext sets a custom context for the Job.
// This context is used for controlling the execution of the job's function.
func (j *Job) WithContext(ctx context.Context) *Job {
	// locking in case you change on the fly but would not recommend
	j.mutex.Lock()
	defer j.mutex.Unlock()
	// cancel the previous context
	if j.cancelFunc != nil {
		j.cancelFunc()
	}
	j.Ctx, j.cancelFunc = context.WithCancel(ctx)
	return j
}

// SetTimezone sets the timezone in which the Job's schedule will be interpreted.
func (j *Job) SetTimezone(loc *time.Location) *Job {
	// locking in case you change on the fly but would not recommend
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.Timezone = loc
	return j
}

// Execute sets the function (Fn) to be executed by the Job.
// The provided function should accept a context.Context parameter.
func (j *Job) Execute(fn func(ctx context.Context)) *Job {
	// locking in case you change on the fly but would not recommend
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.Fn = fn
	return j
}

// Start initiates the execution of the Job according to its schedule.
// The job runs either synchronously or asynchronously based on its Blocking setting.
func (j *Job) Start() {
	j.mutex.Lock()
	if j.Fn == nil || j.isRunning {
		j.mutex.Unlock()
		return
	}
	j.isRunning = true
	j.mutex.Unlock()

	done := j.Ctx.Done()
	go func() {
		for {
			j.mutex.RLock()
			previousRun := j.now()
			// Schedule has a next function that tells you when to run the job next
			// https://pkg.go.dev/github.com/robfig/cron#Schedule
			currentRun := j.Schedule.Next(previousRun)
			timer := time.NewTimer(currentRun.Sub(j.now()))
			isBlocking := j.Blocking
			j.mutex.RUnlock()
			select {
			case <-timer.C:
				if isBlocking {
					j.Fn(j.Ctx)
				} else {
					go j.Fn(j.Ctx)
				}
			case <-done:
				timer.Stop()
				return
			}
		}
	}()
}

// Stop halts the execution of the Job.
// It cancels the Job's context, effectively stopping the running task.
func (j *Job) Stop() {
	j.mutex.Lock()
	if j.isRunning {
		j.isRunning = false
		j.cancelFunc()
	}
	j.mutex.Unlock()
}
