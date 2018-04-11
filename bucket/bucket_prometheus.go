package bucket

import (
	"strings"

	"github.com/fiam/gounidecode/unidecode"
)

// Prometheus struct in an implementation of Bucket interface that produces metric names with labels for given section and operation
type Prometheus struct {
	section   string
	operation string
	success   bool
}

// NewPrometheus builds and returns new Prometheus instance
func NewPrometheus(section string, operation MetricOperation, success, uniDecode bool) *Prometheus {
	operationSanitized := make([]string, cap(operation.operations))
	for k, v := range operation.operations {
		operationSanitized[k] = sanitizeMetricName(v, uniDecode)

		// prometheus doesn't allow _ in then end of metric name
		if operationSanitized[k] == "" {
			operationSanitized = append(operationSanitized[:k], operationSanitized[k+1:]...)
		}
	}

	return &Prometheus{
		section:   sanitizeMetricName(section, uniDecode),
		operation: strings.Join(operationSanitized, "_"),
		success:   success,
	}
}

// Metric builds simple metric name in the form:
//  <section>_<operation-0>_<operation-1>_<operation-2>
func (b *Prometheus) Metric() string {
	return b.section + "_" + b.operation
}

// MetricWithSuffix builds metric name with success suffix in the form:
//  <section>-ok|fail.<operation-0>.<operation-1>.<operation-2>
func (b *Prometheus) MetricWithSuffix() string {
	return b.section + "-" + operationsStatus[b.success] + "." + b.operation
}

// MetricTotal builds simple total metric name in the form:
//  total.<section>
func (b *Prometheus) MetricTotal() string {
	return totalBucket + "_" + b.section
}

// MetricTotalWithSuffix builds total metric name with success suffix in the form
//  total-ok|fail_<section>
func (b *Prometheus) MetricTotalWithSuffix() string {
	return totalBucket + "_" + b.section + "-" + operationsStatus[b.success]
}

func sanitizeMetricName(metric string, uniDecode bool) string {
	if metric == "" {
		return ""
	}

	if uniDecode {
		asciiMetric := unidecode.Unidecode(metric)
		if asciiMetric != metric {
			metric = prefixUnicode + asciiMetric
		}
	}

	metric = strings.Replace(metric, "-", "", -1)

	return strings.Replace(
		// Remove underscores
		strings.Replace(metric, "_", "", -1),
		// and replace dots with single underscore
		".",
		"_",
		-1,
	)
}
