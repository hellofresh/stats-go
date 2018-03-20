package bucket

import (
	"strings"
)

// MetricOperation is a list of metric operations to use for metric
type MetricOperation struct {
	operations  [3]string
	LabelNames  []string
	LabelValues []string
}

// NewMetricOperation  builds and returns new MetricOperation instance with defined label keys
func NewMetricOperation(operations [3]string, labelNames []string) MetricOperation {
	labels := []string{"success"}
	for _, v := range labelNames {
		labels = append(labels, v)
	}

	return MetricOperation{operations: operations, LabelNames: labels}
}

// WithLabels adds label value to existing MetricOperation instance
func (m *MetricOperation) WithLabels(labels ...string) MetricOperation {
	m.LabelValues = labels
	return *m
}

// Plain struct in an implementation of Bucket interface that produces metric names for given section and operation
type Plain struct {
	section   string
	operation string
	success   bool
}

// NewPlain builds and returns new Plain instance
func NewPlain(section string, operation MetricOperation, success, uniDecode bool) *Plain {
	operationSanitized := make([]string, cap(operation.operations))
	for k, v := range operation.operations {
		operationSanitized[k] = SanitizeMetricName(v, uniDecode)
	}
	return &Plain{SanitizeMetricName(section, uniDecode), strings.Join(operationSanitized, "."), success}
}

// Metric builds simple metric name in the form:
//  <section>.<operation-0>.<operation-1>.<operation-2>
func (b *Plain) Metric() string {
	return b.section + "." + b.operation
}

// MetricWithSuffix builds metric name with success suffix in the form:
//  <section>-ok|fail.<operation-0>.<operation-1>.<operation-2>
func (b *Plain) MetricWithSuffix() string {
	return b.section + "-" + operationsStatus[b.success] + "." + b.operation
}

// MetricTotal builds simple total metric name in the form:
//  total.<section>
func (b *Plain) MetricTotal() string {
	return totalBucket + "." + b.section
}

// MetricTotalWithSuffix builds total metric name with success suffix in the form
//  total-ok|fail.<section>
func (b *Plain) MetricTotalWithSuffix() string {
	return totalBucket + "." + b.section + "-" + operationsStatus[b.success]
}
