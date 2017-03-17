package stats

import (
	"net/http"
	"strings"
)

// HttpMetricNameAlterCallback is a type for HTTP Request metric alter handler
type HttpMetricNameAlterCallback func(metricParts MetricOperation, r *http.Request) MetricOperation

// BucketHttpRequest struct in an implementation of Bucket interface that produces metric names for HTTP Request.
// Metrics has the following formats for methods:
//
// * "Metric()" - "<section>.<method>.<path-level-0>.<path-level-1>"
// * "MetricWithSuffix()" - "<section>-ok|fail.<method>.<path-level-0>.<path-level-1>"
// * "TotalRequests()" builds simple total metric name in the form "total.<section>"
// * "MetricTotalWithSuffix()" - builds total metric name with success suffix in the form "total-ok|fail.<section>"
//
// Normally "<section>" is set to "request", but you can use any string value here.
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
