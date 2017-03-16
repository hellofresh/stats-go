package stats

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// LogTimeTracker struct is TimeTracker interface implementation that writes all timings to log
type LogTimeTracker struct {
	timerStart time.Time
}

// Start starts timer
func (t *LogTimeTracker) Start() TimeTracker {
	t.timerStart = time.Now()
	return t
}

// Finish writes elapsed time for metric to log
func (t *LogTimeTracker) Finish(bucket string) {
	log.WithFields(log.Fields{
		"bucket":   bucket,
		"elapsed":  int(time.Now().Sub(t.timerStart) / time.Millisecond),
		"sampling": "ms",
	}).Debug("Muted stats timer send")
}
