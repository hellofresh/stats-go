package stats

import "time"

// TimerMetric is a type for storing single duration metric
type TimerMetric struct {
	Bucket  string
	Elapsed time.Duration
}

// MemoryTimeTracker struct is TimeTracker interface implementation that stores results in memory for further usage
type MemoryTimeTracker struct {
	timerStart time.Time

	bucket  string
	elapsed time.Duration
}

// Start starts timer
func (t *MemoryTimeTracker) Start() TimeTracker {
	t.timerStart = time.Now()
	return t
}

// Finish stores elapsed duration in memory
func (t *MemoryTimeTracker) Finish(bucket string) {
	t.bucket = bucket
	t.elapsed = time.Now().Sub(t.timerStart)
}

// Elapsed returns elapsed duration
func (t *MemoryTimeTracker) Elapsed() TimerMetric {
	return TimerMetric{t.bucket, t.elapsed}
}
