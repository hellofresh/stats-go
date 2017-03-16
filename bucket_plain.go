package stats

import (
	"fmt"
	"strings"
)

type MetricOperation [3]string

type BucketPlain struct {
	section   string
	operation string
	success   bool
}

func NewBucketPlain(section string, operation MetricOperation, success bool) *BucketPlain {
	operationSanitized := make([]string, cap(operation))
	for k, v := range operation {
		operationSanitized[k] = SanitizeMetricName(v)
	}
	return &BucketPlain{SanitizeMetricName(section), strings.Join(operationSanitized, "."), success}
}

func (b *BucketPlain) Metric() string {
	return fmt.Sprintf("%s.%s", b.section, b.operation)
}

func (b *BucketPlain) MetricWithSuffix() string {
	return fmt.Sprintf("%s-%s.%s", b.section, getOperationStatus(b.success), b.operation)
}

func (b *BucketPlain) MetricTotal() string {
	return fmt.Sprintf("%s.%s", totalBucket, b.section)
}

func (b *BucketPlain) MetricTotalWithSuffix() string {
	return fmt.Sprintf("%s.%s-%s", totalBucket, b.section, getOperationStatus(b.success))
}
