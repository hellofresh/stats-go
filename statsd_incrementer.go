package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// StatsdIncrementer struct is Incrementer interface implementation that writes all metrics to statsd
type StatsdIncrementer struct {
	c *statsd.Client
}

// Increment increments metric in statsd
func (i *StatsdIncrementer) Increment(metric string) {
	i.c.Increment(metric)
}

// IncrementN increments metric by n in statsd
func (i *StatsdIncrementer) IncrementN(metric string, n int) {
	i.c.Count(metric, n)
}

// IncrementAll increments all metrics for given bucket in statsd
func (i *StatsdIncrementer) IncrementAll(b Bucket) {
	incrementAll(i, b)
}

// IncrementAllN increments all metrics for given bucket in statsd
func (i *StatsdIncrementer) IncrementAllN(b Bucket, n int) {
	incrementAllN(i, b, n)
}
