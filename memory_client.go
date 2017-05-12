package stats

import (
	"net/http"
	"sync"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/incrementer"
	"github.com/hellofresh/stats-go/state"
	"github.com/hellofresh/stats-go/timer"
)

// MemoryClient is Client implementation for tests
type MemoryClient struct {
	sync.Mutex
	httpMetricCallback bucket.HTTPMetricNameAlterCallback
	httpRequestSection string

	TimerMetrics []timer.Metric
	CountMetrics map[string]int
	StateMetrics map[string]int
}

// NewMemoryClient builds and returns new MemoryClient instance
func NewMemoryClient() *MemoryClient {
	client := &MemoryClient{}
	client.ResetHTTPRequestSection()
	client.resetMetrics()

	return client
}

func (sc *MemoryClient) resetMetrics() {
	sc.TimerMetrics = []timer.Metric{}
	sc.CountMetrics = map[string]int{}
	sc.StateMetrics = map[string]int{}
}

// BuildTimer builds timer to track metric timings
func (sc *MemoryClient) BuildTimer() timer.Timer {
	return &timer.Memory{}
}

// Close resets all collected stats
func (sc *MemoryClient) Close() error {
	sc.resetMetrics()
	return nil
}

// TrackRequest tracks HTTP Request stats
func (sc *MemoryClient) TrackRequest(r *http.Request, t timer.Timer, success bool) Client {
	b := bucket.NewHTTPRequest(sc.httpRequestSection, r, success, sc.httpMetricCallback)
	i := incrementer.NewMemory()

	t.Finish(b.Metric())
	if memoryTimer, ok := t.(*timer.Memory); ok {
		sc.TimerMetrics = append(sc.TimerMetrics, memoryTimer.Elapsed())
	}

	i.IncrementAll(b)
	for metric, value := range i.Metrics() {
		sc.CountMetrics[metric] += value
	}

	return sc
}

// TrackOperation tracks custom operation
func (sc *MemoryClient) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	b := bucket.NewPlain(section, operation, success)
	i := incrementer.NewMemory()

	if nil != t {
		t.Finish(b.MetricWithSuffix())
		if memoryTimer, ok := t.(*timer.Memory); ok {
			sc.TimerMetrics = append(sc.TimerMetrics, memoryTimer.Elapsed())
		}
	}

	i.IncrementAll(b)
	for metric, value := range i.Metrics() {
		sc.CountMetrics[metric] += value
	}

	return sc
}

// TrackOperationN tracks custom operation with n diff
func (sc *MemoryClient) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	b := bucket.NewPlain(section, operation, success)
	i := incrementer.NewMemory()

	if nil != t {
		t.Finish(b.MetricWithSuffix())
		if memoryTimer, ok := t.(*timer.Memory); ok {
			sc.TimerMetrics = append(sc.TimerMetrics, memoryTimer.Elapsed())
		}
	}

	i.IncrementAllN(b, n)
	for metric, value := range i.Metrics() {
		sc.CountMetrics[metric] += value
	}

	return sc
}

// TrackState tracks metric absolute value
func (sc *MemoryClient) TrackState(section string, operation bucket.MetricOperation, value int) Client {
	b := bucket.NewPlain(section, operation, true)
	s := state.NewMemory()

	s.Set(b.Metric(), value)
	for metric, value := range s.Metrics() {
		sc.StateMetrics[metric] = value
	}

	return sc
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (sc *MemoryClient) SetHTTPMetricCallback(callback bucket.HTTPMetricNameAlterCallback) Client {
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
	return sc.SetHTTPRequestSection(bucket.SectionRequest)
}
