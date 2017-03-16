package stats

import "strings"

const (
	totalBucket = "total"

	sectionRequest = "request"

	suffixStatusOk   = "ok"
	suffixStatusFail = "fail"

	MetricEmptyPlaceholder = "-"
	MetricIDPlaceholder    = "-id-"
)

// SanitizeMetricName modifies metric name to work well with statsd
func SanitizeMetricName(metric string) string {
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
