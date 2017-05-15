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

func (c *MemoryClient) resetMetrics() {
	c.TimerMetrics = []timer.Metric{}
	c.CountMetrics = map[string]int{}
	c.StateMetrics = map[string]int{}
}

// BuildTimer builds timer to track metric timings
func (c *MemoryClient) BuildTimer() timer.Timer {
	return &timer.Memory{}
}

// Close resets all collected stats
func (c *MemoryClient) Close() error {
	c.resetMetrics()
	return nil
}

// TrackRequest tracks HTTP Request stats
func (c *MemoryClient) TrackRequest(r *http.Request, t timer.Timer, success bool) Client {
	b := bucket.NewHTTPRequest(c.httpRequestSection, r, success, c.httpMetricCallback)
	i := incrementer.NewMemory()

	t.Finish(b.Metric())
	if memoryTimer, ok := t.(*timer.Memory); ok {
		c.TimerMetrics = append(c.TimerMetrics, memoryTimer.Elapsed())
	}

	i.IncrementAll(b)
	for metric, value := range i.Metrics() {
		c.CountMetrics[metric] += value
	}

	return c
}

// TrackOperation tracks custom operation
func (c *MemoryClient) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	b := bucket.NewPlain(section, operation, success)
	i := incrementer.NewMemory()

	if nil != t {
		t.Finish(b.MetricWithSuffix())
		if memoryTimer, ok := t.(*timer.Memory); ok {
			c.TimerMetrics = append(c.TimerMetrics, memoryTimer.Elapsed())
		}
	}

	i.IncrementAll(b)
	for metric, value := range i.Metrics() {
		c.CountMetrics[metric] += value
	}

	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *MemoryClient) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	b := bucket.NewPlain(section, operation, success)
	i := incrementer.NewMemory()

	if nil != t {
		t.Finish(b.MetricWithSuffix())
		if memoryTimer, ok := t.(*timer.Memory); ok {
			c.TimerMetrics = append(c.TimerMetrics, memoryTimer.Elapsed())
		}
	}

	i.IncrementAllN(b, n)
	for metric, value := range i.Metrics() {
		c.CountMetrics[metric] += value
	}

	return c
}

// TrackState tracks metric absolute value
func (c *MemoryClient) TrackState(section string, operation bucket.MetricOperation, value int) Client {
	b := bucket.NewPlain(section, operation, true)
	s := state.NewMemory()

	s.Set(b.Metric(), value)
	for metric, value := range s.Metrics() {
		c.StateMetrics[metric] = value
	}

	return c
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (c *MemoryClient) SetHTTPMetricCallback(callback bucket.HTTPMetricNameAlterCallback) Client {
	c.Lock()
	defer c.Unlock()

	c.httpMetricCallback = callback
	return c
}

// GetHTTPMetricCallback gets callback handler that allows metric operation alteration for HTTP Request
func (c *MemoryClient) GetHTTPMetricCallback() bucket.HTTPMetricNameAlterCallback {
	c.Lock()
	defer c.Unlock()

	return c.httpMetricCallback
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (c *MemoryClient) SetHTTPRequestSection(section string) Client {
	c.Lock()
	defer c.Unlock()

	c.httpRequestSection = section
	return c
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (c *MemoryClient) ResetHTTPRequestSection() Client {
	return c.SetHTTPRequestSection(bucket.SectionRequest)
}
