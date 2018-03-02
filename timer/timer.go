package timer

import "time"

// Timer is a metric time tracking interface
type Timer interface {
	// Start starts timer
	Start() Timer
	// StartAt starts timer at a given time
	StartAt(time.Time) Timer
	// Finish writes elapsed time for metric
	Finish(bucket string)
}
