package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// TimeTracker is a metric time tracking interface
type TimeTracker interface {
	// Start starts timer
	Start() TimeTracker
	// Finish writes elapsed time for metric
	Finish(bucket string)
}

// NewTimeTracker builds and returns new TimeTracker instance
func NewTimeTracker(c *statsd.Client, muted bool) TimeTracker {
	if muted {
		return &LogTimeTracker{}
	} else {
		return &StatsdTimeTracker{c: c}
	}
}
