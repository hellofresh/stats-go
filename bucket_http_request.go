package stats

import (
	"net/http"
	"strings"
)

// HTTPMetricNameAlterCallback is a type for HTTP Request metric alter handler
type HTTPMetricNameAlterCallback func(metricParts MetricOperation, r *http.Request) MetricOperation

// BucketHTTPRequest struct in an implementation of Bucket interface that produces metric names for HTTP Request.
// Metrics has the following formats for methods:
//  Metric() -> <section>.<method>.<path-level-0>.<path-level-1>
//  MetricWithSuffix() -> <section>-ok|fail.<method>.<path-level-0>.<path-level-1>
//  TotalRequests() -> total.<section>
//  MetricTotalWithSuffix() -> total-ok|fail.<section>
//
// Normally "<section>" is set to "request", but you can use any string value here.
type BucketHTTPRequest struct {
	*BucketPlain
	r        *http.Request
	callback HTTPMetricNameAlterCallback
}

// NewBucketHTTPRequest builds and returns new BucketHTTPRequest instance
func NewBucketHTTPRequest(section string, r *http.Request, success bool, callback HTTPMetricNameAlterCallback) *BucketHTTPRequest {
	operations := getRequestOperations(r, callback)
	return &BucketHTTPRequest{NewBucketPlain(section, operations, success), r, callback}
}

func getRequestOperations(r *http.Request, callback HTTPMetricNameAlterCallback) MetricOperation {
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
