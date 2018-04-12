package incrementer

import (
	"testing"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

type CounterMock struct {
	prometheus.Metric
	prometheus.Collector
}

func (c *CounterMock) Set(float64) {

}

func (c *CounterMock) Inc() {

}

func (c *CounterMock) Add(float64) {

}

type CounterVecMock struct {
	withLabelValuesCalls int
	values               []string
}

func (m *CounterVecMock) GetMetricWithLabelValues(lvs ...string) (prometheus.Counter, error) {
	return nil, nil
}

func (m *CounterVecMock) GetMetricWith(labels prometheus.Labels) (prometheus.Counter, error) {
	return nil, nil
}

func (m *CounterVecMock) WithLabelValues(lvs ...string) prometheus.Counter {
	m.withLabelValuesCalls++
	m.values = lvs
	return &CounterMock{}
}

func (m *CounterVecMock) With(labels prometheus.Labels) prometheus.Counter {
	return nil
}

type CounterFactoryMock struct {
	mock CounterVecMock
}

func (m *CounterFactoryMock) Create(metric string, labelKeys []string) CounterVec {
	m.mock = CounterVecMock{values: labelKeys}
	return &m.mock
}

func TestPrometheus_Increment(t *testing.T) {
	b := bucket.NewPrometheus("section", bucket.NewMetricOperation("o1", "o2", "o3"), true, true)
	m := &CounterFactoryMock{}
	i := NewPrometheus(m)

	i.Increment(b.Metric())
	assert.Equal(t, 1, m.mock.withLabelValuesCalls)
	assert.Equal(t, 0, len(m.mock.values))
}

func TestPrometheus_IncrementWithLabels(t *testing.T) {
	b := bucket.NewPrometheus("section", bucket.NewMetricOperation("o1", "o2", "o3"), true, true)
	m := &CounterFactoryMock{}
	i := NewPrometheus(m)

	i.Increment(b.Metric(), map[string]string{"key1": "value1", "key2": "value2"})
	assert.Equal(t, 1, m.mock.withLabelValuesCalls)
	assert.Equal(t, 2, len(m.mock.values))
}

func TestPrometheus_IncrementN(t *testing.T) {
	b := bucket.NewPrometheus("section", bucket.NewMetricOperation("o1", "o2", "o3"), true, true)
	m := &CounterFactoryMock{}
	i := NewPrometheus(m)

	i.IncrementN(b.Metric(), 123)
	assert.Equal(t, 1, m.mock.withLabelValuesCalls)
	assert.Equal(t, 0, len(m.mock.values))
}

func TestPrometheus_IncrementNWithLabels(t *testing.T) {
	b := bucket.NewPrometheus("section", bucket.NewMetricOperation("o1", "o2", "o3"), true, true)
	m := &CounterFactoryMock{}
	i := NewPrometheus(m)

	i.Increment(b.Metric(), map[string]string{"key1": "value1", "key2": "value2"})
	assert.Equal(t, 1, m.mock.withLabelValuesCalls)
	assert.Equal(t, 2, len(m.mock.values))
}
