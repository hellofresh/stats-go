package stats

import (
	"net/http"
	"strings"
)

type HttpMetricNameAlterCallback func(metricParts MetricOperation, r *http.Request) MetricOperation

type BucketHttpRequest struct {
	*BucketPlain
	r        *http.Request
	callback HttpMetricNameAlterCallback
}

func NewBucketHttpRequest(r *http.Request, success bool, callback HttpMetricNameAlterCallback) *BucketHttpRequest {
	operations := getRequestOperations(r, callback)
	return &BucketHttpRequest{NewBucketPlain("", operations, success), r, callback}
}

func (b *BucketHttpRequest) Request() string {
	b.BucketPlain.section = sectionRequest
	return b.Metric()
}

func (b *BucketHttpRequest) RequestsWithSuffix() string {
	b.BucketPlain.section = sectionRequest
	return b.MetricWithSuffix()
}

func (b *BucketHttpRequest) TotalRequests() string {
	b.BucketPlain.section = sectionRequest
	return b.MetricTotal()
}

func (b *BucketHttpRequest) TotalRequestsWithSuffix() string {
	b.BucketPlain.section = sectionRequest
	return b.MetricTotalWithSuffix()
}

func (b *BucketHttpRequest) RoundTrip() string {
	b.BucketPlain.section = sectionRound
	return b.Metric()
}

func (b *BucketHttpRequest) RoundTripWithSuffix() string {
	b.BucketPlain.section = sectionRound
	return b.MetricWithSuffix()
}

func (b *BucketHttpRequest) TotalRoundTrip() string {
	b.BucketPlain.section = sectionRound
	return b.MetricTotal()
}

func (b *BucketHttpRequest) TotalRoundTripWithSuffix() string {
	b.BucketPlain.section = sectionRound
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
