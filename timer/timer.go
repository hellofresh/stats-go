package timer

// Timer is a metric time tracking interface
type Timer interface {
	// Start starts timer
	Start() Timer
	// Finish writes elapsed time for metric
	Finish(bucket string)
}
