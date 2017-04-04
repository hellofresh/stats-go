package stats

import "net/http"

// Client is an interface for different methods of gathering stats
type Client interface {
	// BuildTimeTracker builds timer to track metric timings
	BuildTimeTracker() TimeTracker
	// Close closes underlying client connection if any
	Close() error

	// TrackRequest tracks HTTP Request stats
	TrackRequest(r *http.Request, tt TimeTracker, success bool) Client

	// TrackOperation tracks custom operation
	TrackOperation(section string, operation MetricOperation, tt TimeTracker, success bool) Client
	// TrackOperation tracks custom operation with n diff
	TrackOperationN(section string, operation MetricOperation, tt TimeTracker, n int, success bool) Client

	// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
	SetHTTPMetricCallback(callback HTTPMetricNameAlterCallback) Client

	// SetHTTPRequestSection sets metric section for HTTP Request metrics
	SetHTTPRequestSection(section string) Client

	// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
	ResetHTTPRequestSection() Client
}
