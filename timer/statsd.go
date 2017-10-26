package timer

import (
	"gopkg.in/alexcesaro/statsd.v2"
)

// StatsD struct is Timer interface implementation that writes all timings to statsd
type StatsD struct {
	timer statsd.Timing
	c     *statsd.Client
}

// NewStatsD creates new statsd timer instance
func NewStatsD(c *statsd.Client) *StatsD {
	return &StatsD{c: c}
}

// Start starts timer
func (t *StatsD) Start() Timer {
	t.timer = t.c.NewTiming()
	return t
}

// Finish writes elapsed time for metric to statsd
func (t *StatsD) Finish(bucket string) {
	t.timer.Send(bucket)
}
