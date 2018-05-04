package incrementer

import (
	"testing"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/stretchr/testify/assert"
)

type metricN struct {
	Metric string
	N      int
}

type mockIncrementer struct {
	callsIncrement  int
	paramsIncrement []string

	callsIncrementN  int
	paramsIncrementN []metricN
}

func (i *mockIncrementer) Increment(metric string, labels ...map[string]string) {
	i.callsIncrement++
	i.paramsIncrement = append(i.paramsIncrement, metric)
}

func (i *mockIncrementer) IncrementN(metric string, n int, labels ...map[string]string) {
	i.callsIncrementN++
	i.paramsIncrementN = append(i.paramsIncrementN, metricN{metric, n})
}

func (i *mockIncrementer) IncrementAll(b bucket.Bucket) {
	incrementAll(i, b)
}

func (i *mockIncrementer) IncrementAllN(b bucket.Bucket, n int) {
	incrementAllN(i, b, n)
}

func Test_incrementAll(t *testing.T) {
	b := bucket.NewPlain("section", bucket.NewMetricOperation("o1", "o2", "o3"), true, true)

	i := &mockIncrementer{}
	i.IncrementAll(b)

	assert.Equal(t, 4, i.callsIncrement)
	assert.Equal(t, b.Metric(), i.paramsIncrement[0])
	assert.Equal(t, b.MetricWithSuffix(), i.paramsIncrement[1])
	assert.Equal(t, b.MetricTotal(), i.paramsIncrement[2])
	assert.Equal(t, b.MetricTotalWithSuffix(), i.paramsIncrement[3])
}

func Test_incrementAllN(t *testing.T) {
	b := bucket.NewPlain("section", bucket.NewMetricOperation("o1", "o2", "o3"), true, true)
	n := 42

	i := &mockIncrementer{}
	i.IncrementAllN(b, n)

	assert.Equal(t, 4, i.callsIncrementN)
	assert.Equal(t, b.Metric(), i.paramsIncrementN[0].Metric)
	assert.Equal(t, b.MetricWithSuffix(), i.paramsIncrementN[1].Metric)
	assert.Equal(t, b.MetricTotal(), i.paramsIncrementN[2].Metric)
	assert.Equal(t, b.MetricTotalWithSuffix(), i.paramsIncrementN[3].Metric)
	assert.Equal(t, n, i.paramsIncrementN[0].N)
	assert.Equal(t, n, i.paramsIncrementN[1].N)
	assert.Equal(t, n, i.paramsIncrementN[2].N)
	assert.Equal(t, n, i.paramsIncrementN[3].N)
}
