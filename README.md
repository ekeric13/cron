# Cron: A Simple Go Cron Scheduling Package

Cron is a Go package for job scheduling which allows you to execute Go functions at specified times or intervals. Inspired by traditional UNIX cron, it extends the functionality with support for more granular time units like seconds and milliseconds, while maintaining simplicity and ease of use.

## Features

- **Simple API**: If you want something more extensive I definitely recommend using [gocron](https://github.com/go-co-op/gocron). The code is essentially all in `./pkg/cron/cron.go` and it is less than 200 lines so you know exactly what you are getting.
- **Flexible Scheduling**: Supports traditional UNIX cron format with extended support for seconds and milliseconds. Uses the cron parser defined in [robfig/cron](https://pkg.go.dev/github.com/robfig/cron?utm_source=godoc#hdr-CRON_Expression_Format)
- **Timezone Awareness**: Schedule jobs in different timezones.
- **Context Support**: Integrates with Go's `context.Context` for job cancellation and timeouts.
- **Blocking/Non-Blocking Execution**: Choose between blocking and non-blocking job execution.
- **Thread-Safe Job Modifications**: Safely modify job settings even after scheduling.

## Getting Started

### Installation

To install the package, use the following command:

```bash
go get github.com/ekeric13/cron
```

### Running Tests

Run the tests using:

```bash
go test ./pkg/...
```

### Running Examples

To see the package in action, run the provided examples:

```bash
go run examples/basic_job/main.go
go run examples/context_job/main.go
go run examples/complex_schedules/main.go
```

## Usage

Here's a quick example to get you started:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/ekeric13/cron/pkg/cron"
)

func main() {
    job := cron.Schedule("* * * * * *").Execute(func(ctx context.Context) {
        fmt.Println("Job executed every second")
    })

    job.Start()

    // Let the job run for 10 seconds
    time.Sleep(10 * time.Second)

    job.Stop()
}
```

## Contributing

Contributions to improve the package are welcome. Please adhere to the following guidelines:

- Write tests for new features and bug fixes.
- Follow the existing coding style and conventions.
- Create a pull request with a clear description of your changes.

That said I would prefer if you just fork it and make changes yourself. And better yet just copy and paste `./pkg/cron/cron.go`, the only dependency is `github.com/robfig/cron/v3`. Personally I much prefer to use libraries that are very easy and quick to grok.

