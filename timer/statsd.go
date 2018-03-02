package timer

import (
	"time"

	"gopkg.in/alexcesaro/statsd.v2"
)

// StatsD struct is Timer interface implementation that writes all timings to statsd
type StatsD struct {
	timerStart time.Time
	c          *statsd.Client
}

// NewStatsD creates new statsd timer instance
func NewStatsD(c *statsd.Client) *StatsD {
	return &StatsD{c: c}
}

// StartAt starts timer at a given time
func (t *StatsD) StartAt(s time.Time) Timer {
	t.timerStart = s
	return t
}

// Start starts timer
func (t *StatsD) Start() Timer {
	t.timerStart = time.Now()
	return t
}

// Finish writes elapsed time for metric to statsd
func (t *StatsD) Finish(bucket string) {
	t.c.Timing(bucket, int(time.Now().Sub(t.timerStart)/time.Millisecond))
}
