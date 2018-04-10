package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrometheus_Metric(t *testing.T) {
	dataProvider := []struct {
		Section   string
		Operation MetricOperation
		Success   bool
		Metric    string
	}{
		{"foo", NewMetricOperation("bar", "baz", "qaz"), true, "foo_bar_baz_qaz"},
		{"foo", NewMetricOperation("bar", "baz", MetricEmptyPlaceholder), true, "foo_bar_baz"},
		{"foo", NewMetricOperation("bar", "dot.baz", MetricEmptyPlaceholder), true, "foo_bar_dot_baz"},
		{"foo", NewMetricOperation("bar", "underscore_baz", MetricEmptyPlaceholder), true, "foo_bar_underscorebaz"},
		{"foo.foo", NewMetricOperation("bar", "underscore_baz", MetricEmptyPlaceholder), true, "foo_foo_bar_underscorebaz"},
	}

	for _, data := range dataProvider {
		b := NewPrometheus(data.Section, data.Operation, data.Success, true)
		assert.Equal(t, data.Metric, b.Metric())
	}
}

func BenchmarkNewPrometheus(b *testing.B) {
	operation := NewMetricOperation("bar", "baz", "qaz")
	for n := 0; n < b.N; n++ {
		NewPrometheus("foo", operation, true, false)
	}
}

func BenchmarkPrometheus_Metric(b *testing.B) {
	bucket := NewPrometheus("foo", NewMetricOperation("bar", "baz", "qaz"), true, false)
	for n := 0; n < b.N; n++ {
		bucket.Metric()
	}
}

func BenchmarkPrometheus_MetricTotal(b *testing.B) {
	bucket := NewPrometheus("foo", NewMetricOperation("bar", "baz", "qaz"), true, false)
	for n := 0; n < b.N; n++ {
		bucket.MetricTotal()
	}
}

func TestPrometheus_MetricTotal(t *testing.T) {
	dataProvider := []struct {
		Section   string
		Operation MetricOperation
		Success   bool
		Metric    string
	}{
		{"foo", NewMetricOperation("bar", "baz", "qaz"), true, "total_foo"},
		{"foo", NewMetricOperation("bar", "baz", "qaz"), false, "total_foo"},
	}

	for _, data := range dataProvider {
		b := NewPrometheus(data.Section, data.Operation, data.Success, true)
		assert.Equal(t, data.Metric, b.MetricTotal())
	}
}
