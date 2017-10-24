package incrementer

import (
	"testing"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/stretchr/testify/assert"
)

func TestMemory_Increment(t *testing.T) {
	i := NewMemory()

	metric1 := "metric-1"
	metric2 := "metric-2"
	metricN := 5

	i.Increment(metric1)
	i.IncrementN(metric2, metricN)

	metrics := i.Metrics()
	assert.Equal(t, 1, metrics[metric1])
	assert.Equal(t, metricN, metrics[metric2])

	i.Increment(metric2)
	i.IncrementN(metric1, metricN)

	metrics = i.Metrics()
	assert.Equal(t, 1+metricN, metrics[metric1])
	assert.Equal(t, metricN+1, metrics[metric2])

	for name := range metrics {
		assert.True(t, name == metric1 || name == metric2)
	}
}

func TestMemory_IncrementAll(t *testing.T) {
	i := NewMemory()
	b1 := bucket.NewPlain("section1", bucket.MetricOperation{"o11", "o12", "o13"}, true, true)
	b2 := bucket.NewPlain("section2", bucket.MetricOperation{"o21", "o22", "o23"}, false, true)
	metricN := 5

	allB1 := []string{b1.Metric(), b1.MetricWithSuffix(), b1.MetricTotal(), b1.MetricTotalWithSuffix()}
	allB2 := []string{b2.Metric(), b2.MetricWithSuffix(), b2.MetricTotal(), b2.MetricTotalWithSuffix()}

	i.IncrementAll(b1)
	i.IncrementAllN(b2, metricN)

	metrics := i.Metrics()
	for i := 0; i < len(allB1); i++ {
		assert.Equal(t, 1, metrics[allB1[i]])
		assert.Equal(t, metricN, metrics[allB2[i]])
	}
}
