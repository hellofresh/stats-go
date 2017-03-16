package stats

import (
	"net/http"
	"strings"
)

// HttpMetricNameAlterCallback is a type for HTTP Request metric alter handler
type HttpMetricNameAlterCallback func(metricParts MetricOperation, r *http.Request) MetricOperation

// BucketHttpRequest struct
type BucketHttpRequest struct {
	*BucketPlain
	r        *http.Request
	callback HttpMetricNameAlterCallback
}

// NewBucketHttpRequest builds and returns new BucketHttpRequest instance
func NewBucketHttpRequest(section string, r *http.Request, success bool, callback HttpMetricNameAlterCallback) *BucketHttpRequest {
	operations := getRequestOperations(r, callback)
	return &BucketHttpRequest{NewBucketPlain(section, operations, success), r, callback}
}

// Request builds simple metric name in the form "request.<method>.<path-level-0>.<path-level-1>"
func (b *BucketHttpRequest) Request() string {
	return b.Metric()
}

// RequestsWithSuffix builds metric name with success suffix in the form "request-ok|fail.<method>.<path-level-0>.<path-level-1>"
func (b *BucketHttpRequest) RequestsWithSuffix() string {
	return b.MetricWithSuffix()
}

// TotalRequests builds simple total metric name in the form total.request"
func (b *BucketHttpRequest) TotalRequests() string {
	return b.MetricTotal()
}

// MetricTotalWithSuffix builds total metric name with success suffix in the form total-ok|fail.request"
func (b *BucketHttpRequest) TotalRequestsWithSuffix() string {
	return b.MetricTotalWithSuffix()
}

func getRequestOperations(r *http.Request, callback HttpMetricNameAlterCallback) MetricOperation {
	metricParts := MetricOperation{strings.ToLower(r.Method), MetricEmptyPlaceholder, MetricEmptyPlaceholder}
	if r.URL.Path != "/" {
		partsFilled := 1
		for _, fragment := range strings.Split(r.URL.Path, "/") {
			if fragment == "" {
				continue
			}

			metricParts[partsFilled] = fragment
			partsFilled++
			if partsFilled >= len(metricParts) {
				break
			}
		}
	}

	if callback != nil {
		metricParts = callback(metricParts, r)
	}

	return metricParts
}
