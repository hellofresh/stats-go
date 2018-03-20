package state

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

type GaugeMock struct {
	prometheus.Metric
	prometheus.Collector
	addCalled      bool
	addCalledValue float64
}

func (c *GaugeMock) Set(float64) {

}

func (c *GaugeMock) Inc() {

}

func (c *GaugeMock) Dec() {

}

func (c *GaugeMock) Add(n float64) {
	c.addCalled = true
	c.addCalledValue = n
}

func (c *GaugeMock) Sub(float64) {

}

type GaugeVecMock struct {
	withLabelValuesCalls int
	values               []string
	gaugeMock            GaugeMock
}

func (m *GaugeVecMock) GetMetricWithLabelValues(lvs ...string) (prometheus.Gauge, error) {
	return nil, nil
}

func (m *GaugeVecMock) GetMetricWith(labels prometheus.Labels) (prometheus.Gauge, error) {
	return nil, nil
}

func (m *GaugeVecMock) WithLabelValues(lvs ...string) prometheus.Gauge {
	m.withLabelValuesCalls++
	m.values = lvs
	return &m.gaugeMock
}

func (m *GaugeVecMock) With(labels prometheus.Labels) prometheus.Gauge {
	return nil
}

type GaugeFactoryMock struct {
	mock GaugeVecMock
}

func (m *GaugeFactoryMock) Create(metric string, labelKeys []string) GaugeVec {
	m.mock = GaugeVecMock{values: labelKeys, gaugeMock: GaugeMock{}}
	return &m.mock
}

func TestPrometheus_Set(t *testing.T) {
	metric1 := "metric1"
	metricState1 := 10

	m := &GaugeFactoryMock{}
	s := NewPrometheus(m)

	s.Set(metric1, metricState1)
	assert.Equal(t, 1, m.mock.withLabelValuesCalls)
	assert.Equal(t, 0, len(m.mock.values))
	assert.Equal(t, true, m.mock.gaugeMock.addCalled)
	assert.Equal(t, float64(10), m.mock.gaugeMock.addCalledValue)
}

func TestPrometheus_SetWithLabels(t *testing.T) {
	metric1 := "metric1"
	metricState1 := 10

	m := &GaugeFactoryMock{}
	s := NewPrometheus(m)

	s.SetWithLabels(metric1, metricState1, []string{"key1"}, []string{"value1"})
	assert.Equal(t, 1, m.mock.withLabelValuesCalls)
	assert.Equal(t, 1, len(m.mock.values))
	assert.Equal(t, true, m.mock.gaugeMock.addCalled)
	assert.Equal(t, float64(10), m.mock.gaugeMock.addCalledValue)
}
