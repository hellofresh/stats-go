<p align="center">
  <a href="https://hellofresh.com">
    <img width="120" src="https://www.hellofresh.de/images/hellofresh/press/HelloFresh_Logo.png">
  </a>
</p>

# hellofresh/stats-go

[![Build Status](https://travis-ci.org/hellofresh/stats-go.svg?branch=master)](https://travis-ci.org/hellofresh/stats-go)
[![GoDoc](https://godoc.org/github.com/hellofresh/stats-go?status.svg)](https://godoc.org/github.com/hellofresh/stats-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/hellofresh/stats-go)](https://goreportcard.com/report/github.com/hellofresh/stats-go)

> Generic Stats library written in Go

This is generic stats library that we use in our projects to collects services stats and then create monitoring
dashboards to track activity and problems.

## Key Features

* Two stats backends - `statsd` for production and `log` for development environment
* Fixed metric sections count for all metrics that allows easy monitoring/alerting setup in `graphana`
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
        statsdClient := stats.NewStatsdStatsClient("statsd-host:8125", "my.app.prefix")
        // explicitly use debug log backend for stats
        mutedClient := stats.NewStatsdStatsClient("", "my.app.prefix")

        // get settings from env to determine backend and prefix
        statsClient := stats.NewStatsdStatsClient(os.Getenv("STATS_DSN"), os.Getenv("STATS_PREFIX"))
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
func NewStatsRequest(sc stats.StatsClient) gin.HandlerFunc {
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
        statsClient := stats.NewStatsdStatsClient(os.Getenv("STATS_DSN"), os.Getenv("STATS_PREFIX"))

        router := gin.Default()
        router.Use(middleware.NewStatsRequest(statsClient))

        router.GET("/", func(c *gin.Context) {
                // will produce "<prefix>.get.-.-" metric
                c.JSON(http.StatusOK, "I'm producing stats!")
        })

        router.Run(":8080")
}
```

#### Generalise resources by type and stripping resource ID

In some cases you do not need to collects metrics for all unique requests, but a single metric for requests of the similar type,
e.g. access time to concrete users pages does not matter a lot, but average access time is important.
`hellofresh/stats-go` allows HTTP Request metric modification and supports IDs filtering out of the box, so
you can get generic metric `get.users.-id-` instead of thousands metrics like `get.users.1`, `get.users.13`,
`get.users.42` etc. that may make your `graphite` suffer from overloading.

To use metric generalisation by second level path ID you can pass `stats.HttpMetricNameAlterCallback` instance to
`stats.StatsClient.SetHttpMetricCallback()`. Also there is a shortcut function `stats.NewHasIDAtSecondLevelCallback()`
that generates callback handler for `stats.SectionsTestsMap`, and shortcut function `stats.ParseSectionsTestsMap`,
that generates sections test map from string, so you can get this values from config.
It accepts list of sections with test callback in the following format: `<section>:<test-callback-name>`.
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

You can register your own test callback functions using `stats.RegisterSectionTest()` function
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
        statsClient := stats.NewStatsdStatsClient(os.Getenv("STATS_DSN"), os.Getenv("STATS_PREFIX")).
                SetHttpMetricCallback(stats.NewHasIDAtSecondLevelCallback(sectionsTestsMap))

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
