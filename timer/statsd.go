package timer

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// Statsd struct is Timer interface implementation that writes all timings to statsd
type Statsd struct {
	timer statsd.Timing
	c     *statsd.Client
}

// Start starts timer
func (t *Statsd) Start() Timer {
	t.timer = t.c.NewTiming()
	return t
}

// Finish writes elapsed time for metric to statsd
func (t *Statsd) Finish(bucket string) {
	t.timer.Send(bucket)
}
