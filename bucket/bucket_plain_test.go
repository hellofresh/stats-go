package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlain_Metric(t *testing.T) {
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
		b := NewPlain(data.Section, data.Operation, data.Success, true)
		assert.Equal(t, data.Metric, b.Metric())
	}
}

func BenchmarkNewPlain(b *testing.B) {
	operation := MetricOperation{"bar", "baz", "qaz"}
	for n := 0; n < b.N; n++ {
		NewPlain("foo", operation, true, false)
	}
}

func BenchmarkNewPlain_unicode(b *testing.B) {
	operation := MetricOperation{"bar", "baz", "qaz"}
	for n := 0; n < b.N; n++ {
		NewPlain("foo", operation, true, true)
	}
}

func BenchmarkPlain_Metric(b *testing.B) {
	bucket := NewPlain("foo", MetricOperation{"bar", "baz", "qaz"}, true, false)
	for n := 0; n < b.N; n++ {
		bucket.Metric()
	}
}

func BenchmarkPlain_MetricWithSuffix(b *testing.B) {
	bucket := NewPlain("foo", MetricOperation{"bar", "baz", "qaz"}, true, false)
	for n := 0; n < b.N; n++ {
		bucket.MetricWithSuffix()
	}
}

func BenchmarkPlain_MetricTotal(b *testing.B) {
	bucket := NewPlain("foo", MetricOperation{"bar", "baz", "qaz"}, true, false)
	for n := 0; n < b.N; n++ {
		bucket.MetricTotal()
	}
}

func BenchmarkPlain_MetricTotalWithSuffix(b *testing.B) {
	bucket := NewPlain("foo", MetricOperation{"bar", "baz", "qaz"}, true, false)
	for n := 0; n < b.N; n++ {
		bucket.MetricTotalWithSuffix()
	}
}

func TestPlain_MetricWithSuffix(t *testing.T) {
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
		b := NewPlain(data.Section, data.Operation, data.Success, true)
		assert.Equal(t, data.Metric, b.MetricWithSuffix())
	}
}

func TestPlain_MetricTotal(t *testing.T) {
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
		b := NewPlain(data.Section, data.Operation, data.Success, true)
		assert.Equal(t, data.Metric, b.MetricTotal())
	}
}

func TestPlain_MetricTotalWithSuffix(t *testing.T) {
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
		b := NewPlain(data.Section, data.Operation, data.Success, true)
		assert.Equal(t, data.Metric, b.MetricTotalWithSuffix())
	}
}
