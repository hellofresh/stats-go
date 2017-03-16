package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// StatsdIncrementer struct is Incrementer interface implementation that writes all metrics to statsd
type StatsdIncrementer struct {
	c *statsd.Client
}

// Increment writes given metric to statsd
func (t *StatsdIncrementer) Increment(bucket string) {
	t.c.Increment(bucket)
}
