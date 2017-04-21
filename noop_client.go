package stats

import "net/http"

// NoopClient is Client implementation that does literally nothing
type NoopClient struct{}

// NewNoopClient builds and returns new NoopClient instance
func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

// BuildTimeTracker builds timer to track metric timings
func (c *NoopClient) BuildTimeTracker() TimeTracker {
	return &MemoryTimeTracker{}
}

// Close closes underlying client connection if any
func (c *NoopClient) Close() error {
	return nil
}

// TrackRequest tracks HTTP Request stats
func (c *NoopClient) TrackRequest(r *http.Request, tt TimeTracker, success bool) Client {
	return c
}

// TrackOperation tracks custom operation
func (c *NoopClient) TrackOperation(section string, operation MetricOperation, tt TimeTracker, success bool) Client {
	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *NoopClient) TrackOperationN(section string, operation MetricOperation, tt TimeTracker, n int, success bool) Client {
	return c
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (c *NoopClient) SetHTTPMetricCallback(callback HTTPMetricNameAlterCallback) Client {
	return c
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (c *NoopClient) SetHTTPRequestSection(section string) Client {
	return c
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (c *NoopClient) ResetHTTPRequestSection() Client {
	return c
}
