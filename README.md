<p align="center">
  <a href="https://hellofresh.com">
    <img width="120" src="https://www.hellofresh.de/images/hellofresh/press/HelloFresh_Logo.png">
  </a>
</p>

# hellofresh/stats-go

[![Build Status](https://travis-ci.org/hellofresh/stats-go.svg?branch=master)](https://travis-ci.org/hellofresh/stats-go)
[![Coverage Status](https://coveralls.io/repos/github/hellofresh/stats-go/badge.svg?branch=master)](https://coveralls.io/github/hellofresh/stats-go?branch=master)
[![GoDoc](https://godoc.org/github.com/hellofresh/stats-go?status.svg)](https://godoc.org/github.com/hellofresh/stats-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/hellofresh/stats-go)](https://goreportcard.com/report/github.com/hellofresh/stats-go)

> Generic Stats library written in Go

This is generic stats library that we at HelloFresh use in our projects to collect services' stats and then create monitoring
dashboards to track activity and problems.

## Key Features

* Several stats backends:
  * `log` for development environment
  * `statsd` for production (with fallback to `log` if statsd server is not available)
  * `memory` for testing purpose, to track stats operations in unit tests
  * `noop` for environments that do not require any stats gathering
* Fixed metric sections count for all metrics to allow easy monitoring/alerting setup in `grafana`
* Easy to build HTTP requests metrics - timing and count
* Generalise or modify HTTP Requests metric - e.g. skip ID part

## Installation

```sh
go get -u github.com/hellofresh/stats-go
```

## Usage

#### Instance creation

```go
package main

import (
        "os"

        "github.com/hellofresh/stats-go"
)

func main() {
        // client that tries to connect to statsd service, fallback to debug log backend if fails to connect
        statsdClient, _ := stats.NewClient("statsd://statsd-host:8125", "my.app.prefix")
        defer statsdClient.Close()

        // debug log backend for stats
        logClient, _ := stats.NewClient("log://", "")
        defer logClient.Close()

        // memory backend to track operations in unit tests
        memoryClient, _ := stats.NewClient("memory://", "")
        defer memoryClient.Close()

        // noop backend to ignore all stats
        noopClient, _ := stats.NewClient("noop://", "")
        defer noopClient.Close()

        // client that tries to connect to statsd service, fallback to debug log backend if fails to connect
        // format for backward compatibility with previous version
        legacyStatsdClient, _ := stats.NewClient("statsd-host:8125", "my.app.prefix")
        defer legacyStatsdClient.Close()

        // debug log backend for stats
        // format for backward compatibility with previous version
        legacyLogClient, _ := stats.NewClient("", "")
        defer legacyLogClient.Close()

        // get settings from env to determine backend and prefix
        statsClient, _ := stats.NewClient(os.Getenv("STATS_DSN"), os.Getenv("STATS_PREFIX"))
        defer statsClient.Close()
}
```

#### Count metrics manually

```go
timing := statsClient.BuildTimeTracker().Start()
operations := statsClient.MetricOperation{"orders", "order", "create"}
err := orderService.Create(...)
statsClient.TrackOperation("ordering", operations, timing, err == nil)
```

#### Track requests metrics with middleware, e.g. for [Gin Web Framework](https://github.com/gin-gonic/gin)

```go
package middleware

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	stats "github.com/hellofresh/stats-go"
)

// NewStatsRequest returns a middleware handler function.
func NewStatsRequest(sc stats.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithField("path", c.Request.URL.Path).Debug("Starting Stats middleware")

		timing := sc.BuildTimeTracker().Start()

		c.Next()

		success := c.Writer.Status() < http.StatusBadRequest
		log.WithFields(log.Fields{"request_url": c.Request.URL.Path}).Debug("Track request stats")
		sc.TrackRequest(c.Request, timing, success)
	}
}
```

```go
package main

import (
        "net/http"
        "os"

        "github.com/example/app/middleware"
        "github.com/gin-gonic/gin"
        stats "github.com/hellofresh/stats-go"
)

func main() {
        statsClient := stats.NewStatsdClient(os.Getenv("STATS_DSN"), os.Getenv("STATS_PREFIX"))
        defer statsClient.Close()

        router := gin.Default()
        router.Use(middleware.NewStatsRequest(statsClient))

        router.GET("/", func(c *gin.Context) {
                // will produce "<prefix>.get.-.-" metric
                c.JSON(http.StatusOK, "I'm producing stats!")
        })

        router.Run(":8080")
}
```

#### Usage in unit tests

```go
package foo

import "github.com/hellofresh/stats-go"

const sectionStatsFoo = "foo"

func DoSomeJob(statsClient stats.Client) error {
        tt := statsClient.BuildTimeTracker().Start()
        operation := stats.MetricOperation{"do", "some", "job"}

        result, err := doSomeRealJobHere()
        statsClient.TrackOperation(sectionStatsFoo, operation, tt, result)

        return err
}
```

```go
package foo

import (
        "testing"

        "github.com/hellofresh/stats-go"
        "github.com/stretchr/testify/assert"
)

func TestDoSomeJob(t *testing.T) {
        statsClient, _ := stats.NewClient("memory://", "") 
        
        err := DoSomeJob(statsClient)
        assert.Nil(t, err)
        
        statsMemory, _ := statsClient.(stats.MemoryClient)
        assert.Equal(t, 1, len(statsMemory.TimeMetrics))
        assert.Equal(t, "foo-ok.do.some.job", statsMemory.TimeMetrics[0].Bucket)
        assert.Equal(t, 1, statsMemory.CountMetrics["foo-ok.do.some.job"])
}
```

#### Generalise resources by type and stripping resource ID

In some cases you do not need to collect metrics for all unique requests, but a single metric for requests of the similar type,
e.g. access time to concrete users pages does not matter a lot, but average access time is important.
`hellofresh/stats-go` allows HTTP Request metric modification and supports ID filtering out of the box, so
you can get generic metric `get.users.-id-` instead thousands of metrics like `get.users.1`, `get.users.13`,
`get.users.42` etc. that may make your `graphite` suffer from overloading.

To use metric generalisation by second level path ID, you can pass `stats.HttpMetricNameAlterCallback` instance to
`stats.Client.SetHttpMetricCallback()`. Also there is a shortcut function `stats.NewHasIDAtSecondLevelCallback()`
that generates a callback handler for `stats.SectionsTestsMap`, and shortcut function `stats.ParseSectionsTestsMap`,
that generates sections test map from string, so you can get these values from config.
It accepts a list of sections with test callback in the following format: `<section>:<test-callback-name>`.
You can use either double colon or new line character as section-callback pairs separator, so all of the following
forms are correct:

* `<section-0>:<test-callback-name-0>:<section-1>:<test-callback-name-1>:<section-2>:<test-callback-name-2>`
* `<section-0>:<test-callback-name-0>\n<section-1>:<test-callback-name-1>\n<section-2>:<test-callback-name-2>`
* `<section-0>:<test-callback-name-0>:<section-1>:<test-callback-name-1>\n<section-2>:<test-callback-name-2>`

Currently the following test callbacks are implemented:

* `true` - second path level is always treated as ID,
  e.g. `/users/13` -> `users.-id-`, `/users/search` -> `users.-id-`, `/users` -> `users.-id-`
* `numeric` - only numeric second path level is interpreted as ID,
  e.g. `/users/13` -> `users.-id-`, `/users/search` -> `users.search`
* `not_empty` - only not empty second path level is interpreted as ID,
  e.g. `/users/13` -> `users.-id-`, `/users` -> `users.-`

You can register your own test callback functions using the `stats.RegisterSectionTest()` function
before parsing sections map from string.

```go
package main

import (
        "net/http"
        "os"

        "github.com/example/app/middleware"
        "github.com/gin-gonic/gin"
        stats "github.com/hellofresh/stats-go"
)

func main() {
        // STATS_IDS=users:not_empty:clients:numeric
        sectionsTestsMap, err := stats.ParseSectionsTestsMap(os.Getenv("STATS_IDS"))
        if err != nil {
                sectionsTestsMap = map[stats.PathSection]stats.SectionTestDefinition{}
        }
        statsClient, _ := stats.NewClient(os.Getenv("STATS_DSN"), os.Getenv("STATS_PREFIX"))
        statsClient.SetHTTPMetricCallback(stats.NewHasIDAtSecondLevelCallback(sectionsTestsMap))
        defer statsClient.Close()

        router := gin.Default()
        router.Use(middleware.NewStatsRequest(statsClient))

        router.GET("/users", func(c *gin.Context) {
                // will produce "<prefix>.get.users.-" metric
                c.JSON(http.StatusOK, "Get the userslist")
        })
        router.GET("/users/:id", func(c *gin.Context) {
                // will produce "<prefix>.get.users.-id-" metric 
                c.JSON(http.StatusOK, "Get the user ID " + c.Params.ByName("id"))
        })
        router.GET("/clients/:id", func(c *gin.Context) {
                // will produce "<prefix>.get.clients.-id-" metric
                c.JSON(http.StatusOK, "Get the client ID " + c.Params.ByName("id"))
        })

        router.Run(":8080")
}
```

## Contributing

To start contributing, please check [CONTRIBUTING](CONTRIBUTING.md).

## Documentation

* `hellofresh/stats-go` Docs: https://godoc.org/github.com/hellofresh/stats-go
* Go lang: https://golang.org/
