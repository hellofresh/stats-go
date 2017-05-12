package state

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// State is a metric state interface
type State interface {
	// Set sets metric state
	Set(metric string, n int)
}

// New builds and returns new State instance
func New(c *statsd.Client, muted bool) State {
	if muted {
		return &Log{}
	}
	return &Statsd{c}
}
