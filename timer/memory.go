package timer

import "time"

// Metric is a type for storing single duration metric
type Metric struct {
	Bucket  string
	Elapsed time.Duration
}

// Memory struct is Timer interface implementation that stores results in memory for further usage
type Memory struct {
	timerStart time.Time

	bucket  string
	elapsed time.Duration
}

// Start starts timer
func (t *Memory) Start() Timer {
	t.timerStart = time.Now()
	return t
}

// Finish stores elapsed duration in memory
func (t *Memory) Finish(bucket string) {
	t.bucket = bucket
	t.elapsed = time.Now().Sub(t.timerStart)
}

// Elapsed returns elapsed duration
func (t *Memory) Elapsed() Metric {
	return Metric{t.bucket, t.elapsed}
}