package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// Incrementer is a metric incrementer interface
type Incrementer interface {
	// Increment increments metric for given bucket
	Increment(bucket string)
}

// NewIncrementer builds and returns new Incrementer instance
func NewIncrementer(c *statsd.Client, muted bool) Incrementer {
	if muted {
		return &LogIncrementer{}
	} else {
		return &StatsdIncrementer{c}
	}
}
