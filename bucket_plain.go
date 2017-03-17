package stats

import (
	"fmt"
	"strings"
)

// MetricOperation is a list of metric operations to use for metric
type MetricOperation [3]string

// BucketPlain struct in an implementation of Bucket interface that produces metric names for given section and operation
type BucketPlain struct {
	section   string
	operation string
	success   bool
}

// NewBucketPlain builds and returns new BucketPlain instance
func NewBucketPlain(section string, operation MetricOperation, success bool) *BucketPlain {
	operationSanitized := make([]string, cap(operation))
	for k, v := range operation {
		operationSanitized[k] = SanitizeMetricName(v)
	}
	return &BucketPlain{SanitizeMetricName(section), strings.Join(operationSanitized, "."), success}
}

// Metric builds simple metric name in the form "<section>.<operation-0>.<operation-1>.<operation-2>"
func (b *BucketPlain) Metric() string {
	return fmt.Sprintf("%s.%s", b.section, b.operation)
}

// MetricWithSuffix builds metric name with success suffix in the form "<section>-ok|fail.<operation-0>.<operation-1>.<operation-2>"
func (b *BucketPlain) MetricWithSuffix() string {
	return fmt.Sprintf("%s-%s.%s", b.section, getOperationStatus(b.success), b.operation)
}

// MetricTotal builds simple total metric name in the form total.<section>"
func (b *BucketPlain) MetricTotal() string {
	return fmt.Sprintf("%s.%s", totalBucket, b.section)
}

// MetricTotalWithSuffix builds total metric name with success suffix in the form total-ok|fail.<section>"
func (b *BucketPlain) MetricTotalWithSuffix() string {
	return fmt.Sprintf("%s.%s-%s", totalBucket, b.section, getOperationStatus(b.success))
}
