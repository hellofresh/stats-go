package stats

import "strings"

const (
	totalBucket = "total"

	sectionRequest = "request"

	suffixStatusOk   = "ok"
	suffixStatusFail = "fail"

	// MetricEmptyPlaceholder is a string placeholder for empty (unset) sections of operation
	MetricEmptyPlaceholder = "-"

	// MetricIDPlaceholder is a string placeholder for ID section of operation if any
	MetricIDPlaceholder = "-id-"
)

// Bucket is an interface for building metric names for operations
type Bucket interface {
	// Metric builds simple metric name in the form "<section>.<operation-0>.<operation-1>.<operation-2>"
	Metric() string

	// MetricWithSuffix builds metric name with success suffix in the form "<section>-ok|fail.<operation-0>.<operation-1>.<operation-2>"
	MetricWithSuffix() string

	// MetricTotal builds simple total metric name in the form total.<section>"
	MetricTotal() string

	// MetricTotalWithSuffix builds total metric name with success suffix in the form total-ok|fail.<section>"
	MetricTotalWithSuffix() string
}

// SanitizeMetricName modifies metric name to work well with statsd
func SanitizeMetricName(metric string) string {
	if metric == "" {
		return MetricEmptyPlaceholder
	}

	return strings.Replace(
		// Double underscores
		strings.Replace(metric, "_", "__", -1),
		// and replace dots with single underscore
		".",
		"_",
		-1,
	)
}

func getOperationStatus(success bool) string {
	return map[bool]string{true: suffixStatusOk, false: suffixStatusFail}[success]
}
