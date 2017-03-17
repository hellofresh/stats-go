package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// Incrementer is a metric incrementer interface
type Incrementer interface {
	// Increment increments metric
	Increment(metric string)

	// IncrementN increments metric by n
	IncrementN(metric string, n int)

	// Increment increments all metrics for given bucket
	IncrementAll(b Bucket)

	// Increment increments all metrics for given bucket by n
	IncrementAllN(b Bucket, n int)
}

// NewIncrementer builds and returns new Incrementer instance
func NewIncrementer(c *statsd.Client, muted bool) Incrementer {
	if muted {
		return &LogIncrementer{}
	}
	return &StatsdIncrementer{c}
}
