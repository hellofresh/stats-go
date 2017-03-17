package stats

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBucketPlain_Metric(t *testing.T) {
	dataProvider := []struct {
		Section   string
		Operation MetricOperation
		Success   bool
		Metric    string
	}{
		{"foo", MetricOperation{"bar", "baz", "qaz"}, true, "foo.bar.baz.qaz"},
		{"foo", MetricOperation{"bar", "baz", MetricEmptyPlaceholder}, true, "foo.bar.baz.-"},
		{"foo", MetricOperation{"bar", "dot.baz", MetricEmptyPlaceholder}, true, "foo.bar.dot_baz.-"},
		{"foo", MetricOperation{"bar", "underscore_baz", MetricEmptyPlaceholder}, true, "foo.bar.underscore__baz.-"},
		{"foo.foo", MetricOperation{"bar", "underscore_baz", MetricEmptyPlaceholder}, true, "foo_foo.bar.underscore__baz.-"},
	}

	for _, data := range dataProvider {
		b := NewBucketPlain(data.Section, data.Operation, data.Success)
		assert.Equal(t, data.Metric, b.Metric())
	}
}

func TestBucketPlain_MetricWithSuffix(t *testing.T) {
	dataProvider := []struct {
		Section   string
		Operation MetricOperation
		Success   bool
		Metric    string
	}{
		{"foo", MetricOperation{"bar", "baz", "qaz"}, true, "foo-ok.bar.baz.qaz"},
		{"foo", MetricOperation{"bar", "baz", "qaz"}, false, "foo-fail.bar.baz.qaz"},
	}

	for _, data := range dataProvider {
		b := NewBucketPlain(data.Section, data.Operation, data.Success)
		assert.Equal(t, data.Metric, b.MetricWithSuffix())
	}
}

func TestBucketPlain_MetricTotal(t *testing.T) {
	dataProvider := []struct {
		Section   string
		Operation MetricOperation
		Success   bool
		Metric    string
	}{
		{"foo", MetricOperation{"bar", "baz", "qaz"}, true, "total.foo"},
		{"foo", MetricOperation{"bar", "baz", "qaz"}, false, "total.foo"},
	}

	for _, data := range dataProvider {
		b := NewBucketPlain(data.Section, data.Operation, data.Success)
		assert.Equal(t, data.Metric, b.MetricTotal())
	}
}

func TestBucketPlain_MetricTotalWithSuffix(t *testing.T) {
	dataProvider := []struct {
		Section   string
		Operation MetricOperation
		Success   bool
		Metric    string
	}{
		{"foo", MetricOperation{"bar", "baz", "qaz"}, true, "total.foo-ok"},
		{"foo", MetricOperation{"bar", "baz", "qaz"}, false, "total.foo-fail"},
	}

	for _, data := range dataProvider {
		b := NewBucketPlain(data.Section, data.Operation, data.Success)
		assert.Equal(t, data.Metric, b.MetricTotalWithSuffix())
	}
}
