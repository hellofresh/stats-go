package timer

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// Timer is a metric time tracking interface
type Timer interface {
	// Start starts timer
	Start() Timer
	// Finish writes elapsed time for metric
	Finish(bucket string)
}

// New builds and returns new Timer instance
func New(c *statsd.Client, muted bool) Timer {
	if muted {
		return &Log{}
	}
	return &Statsd{c: c}
}
