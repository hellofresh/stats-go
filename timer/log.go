package timer

import (
	"time"

	"github.com/hellofresh/stats-go/log"
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
	log.Log("Stats timer finished", map[string]interface{}{
		"bucket":  bucket,
		"elapsed": time.Now().Sub(t.timerStart).String(),
	}, nil)
}
