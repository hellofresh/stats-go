package stats

import "strings"

const (
	totalBucket = "total"

	sectionRequest = "request"
	sectionRound   = "round"

	suffixStatusOk   = "ok"
	suffixStatusFail = "fail"

	MetricEmptyPlaceholder = "-"
	MetricIDPlaceholder    = "-id-"
)

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
