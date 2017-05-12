package timer

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// Log struct is Timer interface implementation that writes all timings to log
type Log struct {
	timerStart time.Time
}

// Start starts timer
func (t *Log) Start() Timer {
	t.timerStart = time.Now()
	return t
}

// Finish writes elapsed time for metric to log
func (t *Log) Finish(bucket string) {
	log.WithFields(log.Fields{
		"bucket":  bucket,
		"elapsed": time.Now().Sub(t.timerStart).String(),
	}).Debug("Muted stats timer send")
}
