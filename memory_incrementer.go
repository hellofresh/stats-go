package stats

// MemoryIncrementer struct is Incrementer interface implementation that stores results in memory for further usage
type MemoryIncrementer struct {
	metrics map[string]int
}

// NewMemoryIncrementer builds and returns new MemoryIncrementer instance
func NewMemoryIncrementer() *MemoryIncrementer {
	return &MemoryIncrementer{make(map[string]int)}
}

// Increment increments given metric in memory
func (i *MemoryIncrementer) Increment(metric string) {
	i.metrics[metric]++
}

// IncrementN increments given metric in memory
func (i *MemoryIncrementer) IncrementN(metric string, n int) {
	i.metrics[metric] += n
}

// IncrementAll increments all metrics for given bucket in memory
func (i *MemoryIncrementer) IncrementAll(b Bucket) {
	incrementAll(i, b)
}

// IncrementAllN increments all metrics for given bucket in memory
func (i *MemoryIncrementer) IncrementAllN(b Bucket, n int) {
	incrementAllN(i, b, n)
}

// Metrics returns all previously stored metrics
func (i *MemoryIncrementer) Metrics() map[string]int {
	return i.metrics
}
