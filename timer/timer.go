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
	// FinishWithLabels writes elapsed time for metric
	FinishWithLabels(bucket string, labels map[string]string)
}
