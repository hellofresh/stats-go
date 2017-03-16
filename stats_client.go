package stats

import "net/http"

type StatsClient interface {
	BuildTimeTracker() TimeTracker
	Close()

	TrackRequest(r *http.Request, tt TimeTracker, success bool)
	TrackRoundTrip(r *http.Request, tt TimeTracker, success bool)

	TrackOperation(section string, operation MetricOperation, tt TimeTracker, success bool)

	SetHttpMetricCallback(callback HttpMetricNameAlterCallback) StatsClient
}
