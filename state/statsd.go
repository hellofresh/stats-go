package state

import statsd "gopkg.in/alexcesaro/statsd.v2"

// Statsd struct is State interface implementation that writes all states to statsd gauge
type Statsd struct {
	c *statsd.Client
}

// Set sets metric state
func (s *Statsd) Set(metric string, n int) {
	s.c.Gauge(metric, n)
}
