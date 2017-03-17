package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// StatsdTimeTracker struct is TimeTracker interface implementation that writes all timings to statsd
type StatsdTimeTracker struct {
	timer statsd.Timing
	c     *statsd.Client
}

// Start starts timer
func (t *StatsdTimeTracker) Start() TimeTracker {
	t.timer = t.c.NewTiming()
	return t
}

// Finish writes elapsed time for metric to statsd
func (t *StatsdTimeTracker) Finish(bucket string) {
	t.timer.Send(bucket)
}
