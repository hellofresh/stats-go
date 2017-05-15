package stats

import (
	"net/http"
	"sync"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/timer"
)

// NoopClient is Client implementation that does literally nothing
type NoopClient struct {
	sync.Mutex
	httpMetricCallback bucket.HTTPMetricNameAlterCallback
}

// NewNoopClient builds and returns new NoopClient instance
func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

// BuildTimer builds timer to track metric timings
func (c *NoopClient) BuildTimer() timer.Timer {
	return &timer.Memory{}
}

// Close closes underlying client connection if any
func (c *NoopClient) Close() error {
	return nil
}

// TrackRequest tracks HTTP Request stats
func (c *NoopClient) TrackRequest(r *http.Request, t timer.Timer, success bool) Client {
	return c
}

// TrackOperation tracks custom operation
func (c *NoopClient) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *NoopClient) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	return c
}

// TrackState tracks metric absolute value
func (c *NoopClient) TrackState(section string, operation bucket.MetricOperation, value int) Client {
	return c
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (c *NoopClient) SetHTTPMetricCallback(callback bucket.HTTPMetricNameAlterCallback) Client {
	c.Lock()
	defer c.Unlock()

	c.httpMetricCallback = callback
	return c
}

// GetHTTPMetricCallback gets callback handler that allows metric operation alteration for HTTP Request
func (c *NoopClient) GetHTTPMetricCallback() bucket.HTTPMetricNameAlterCallback {
	c.Lock()
	defer c.Unlock()

	return c.httpMetricCallback
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (c *NoopClient) SetHTTPRequestSection(section string) Client {
	return c
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (c *NoopClient) ResetHTTPRequestSection() Client {
	return c
}
