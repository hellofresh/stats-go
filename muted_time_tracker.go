package stats

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

type MutedTimeTracker struct {
	timerStart time.Time
}

func (t *MutedTimeTracker) Start() TimeTracker {
	t.timerStart = time.Now()
	return t
}

func (t *MutedTimeTracker) Finish(bucket string) {
	log.WithFields(log.Fields{
		"bucket":   bucket,
		"elapsed":  int(time.Now().Sub(t.timerStart) / time.Millisecond),
		"sampling": "ms",
	}).Debug("Muted stats timer send")
}
