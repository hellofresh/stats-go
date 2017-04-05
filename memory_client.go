package stats

import (
	"net/http"
	"sync"
)

// MemoryClient is Client implementation for tests
type MemoryClient struct {
	sync.Mutex
	httpMetricCallback HTTPMetricNameAlterCallback
	httpRequestSection string

	TimeMetrics  []TimerMetric
	CountMetrics map[string]int
}

// NewMemoryClient builds and returns new MemoryClient instance
func NewMemoryClient() *MemoryClient {
	client := &MemoryClient{}
	client.ResetHTTPRequestSection()
	client.resetMetrics()

	return client
}

func (sc *MemoryClient) resetMetrics() {
	sc.TimeMetrics = []TimerMetric{}
	sc.CountMetrics = map[string]int{}
}

// BuildTimeTracker builds timer to track metric timings
func (sc *MemoryClient) BuildTimeTracker() TimeTracker {
	return &MemoryTimeTracker{}
}

// Close resets all collected stats
func (sc *MemoryClient) Close() error {
	sc.resetMetrics()
	return nil
}

// TrackRequest tracks HTTP Request stats
func (sc *MemoryClient) TrackRequest(r *http.Request, tt TimeTracker, success bool) Client {
	b := NewBucketHTTPRequest(sc.httpRequestSection, r, success, sc.httpMetricCallback)
	i := NewMemoryIncrementer()

	tt.Finish(b.Metric())
	if memoryTimer, ok := tt.(*MemoryTimeTracker); ok {
		sc.TimeMetrics = append(sc.TimeMetrics, memoryTimer.Elapsed())
	}

	i.IncrementAll(b)
	for metric, value := range i.Metrics() {
		sc.CountMetrics[metric] += value
	}

	return sc
}

// TrackOperation tracks custom operation
func (sc *MemoryClient) TrackOperation(section string, operation MetricOperation, tt TimeTracker, success bool) Client {
	b := NewBucketPlain(section, operation, success)
	i := NewMemoryIncrementer()

	if nil != tt {
		tt.Finish(b.MetricWithSuffix())
		if memoryTimer, ok := tt.(*MemoryTimeTracker); ok {
			sc.TimeMetrics = append(sc.TimeMetrics, memoryTimer.Elapsed())
		}
	}

	i.IncrementAll(b)
	for metric, value := range i.Metrics() {
		sc.CountMetrics[metric] += value
	}

	return sc
}

// TrackOperationN tracks custom operation with n diff
func (sc *MemoryClient) TrackOperationN(section string, operation MetricOperation, tt TimeTracker, n int, success bool) Client {
	b := NewBucketPlain(section, operation, success)
	i := NewMemoryIncrementer()

	if nil != tt {
		tt.Finish(b.MetricWithSuffix())
		if memoryTimer, ok := tt.(*MemoryTimeTracker); ok {
			sc.TimeMetrics = append(sc.TimeMetrics, memoryTimer.Elapsed())
		}
	}

	i.IncrementAllN(b, n)
	for metric, value := range i.Metrics() {
		sc.CountMetrics[metric] += value
	}

	return sc
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (sc *MemoryClient) SetHTTPMetricCallback(callback HTTPMetricNameAlterCallback) Client {
	sc.Lock()
	defer sc.Unlock()

	sc.httpMetricCallback = callback
	return sc
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (sc *MemoryClient) SetHTTPRequestSection(section string) Client {
	sc.Lock()
	defer sc.Unlock()

	sc.httpRequestSection = section
	return sc
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (sc *MemoryClient) ResetHTTPRequestSection() Client {
	return sc.SetHTTPRequestSection(sectionRequest)
}
