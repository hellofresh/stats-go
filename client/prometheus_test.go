package client

import (
	"testing"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/incrementer"
	"github.com/hellofresh/stats-go/state"
	"github.com/hellofresh/stats-go/timer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

// Mock incrementer object
type mockIncrementer struct {
	incrementMethodCalled  int
	incrementMethodMetrics []string
	incrementMethodLabels  []map[string]string

	incrementNMethodCalled  int
	incrementNMethodMetrics []string
	incrementNMethodLabels  []map[string]string

	incrementALlMethodCalled  int
	incrementALlNMethodCalled int
}

func (i *mockIncrementer) Increment(metric string, labels ...map[string]string) {
	i.incrementMethodCalled++
	i.incrementMethodMetrics = append(i.incrementMethodMetrics, metric)
	if labels[0] != nil {
		i.incrementMethodLabels = append(i.incrementMethodLabels, labels[0])
	}
}

func (i *mockIncrementer) IncrementN(metric string, n int, labels ...map[string]string) {
	i.incrementNMethodCalled++
	i.incrementNMethodMetrics = append(i.incrementNMethodMetrics, metric)
	if labels[0] != nil {
		i.incrementNMethodLabels = append(i.incrementNMethodLabels, labels[0])
	}
}

func (i *mockIncrementer) IncrementAll(b bucket.Bucket) {
	i.incrementALlMethodCalled++
}

func (i *mockIncrementer) IncrementAllN(b bucket.Bucket, n int) {
	i.incrementALlNMethodCalled++
}

// Mock IncrementFactory object
type mockIncrementerFactory struct {
	inc *mockIncrementer

	createMethodCalled int
}

func newMockIncrementerFactory() *mockIncrementerFactory {
	return &mockIncrementerFactory{createMethodCalled: 0}
}

func (m *mockIncrementerFactory) Create() incrementer.Incrementer {
	m.createMethodCalled++
	if m.inc == nil {
		m.inc = &mockIncrementer{
			incrementMethodCalled:  0,
			incrementMethodMetrics: []string{},
			incrementMethodLabels:  []map[string]string{},

			incrementNMethodCalled:  0,
			incrementNMethodMetrics: []string{},
			incrementNMethodLabels:  []map[string]string{},

			incrementALlMethodCalled:  0,
			incrementALlNMethodCalled: 0,
		}
	}
	return m.inc
}

// Mock state object
type mockState struct {
	setMethodCalled  int
	setMethodNumbers []int
	setMethodMetrics []string
	setMethodLabels  []map[string]string
}

func (s *mockState) Set(metric string, n int, labels ...map[string]string) {
	s.setMethodCalled++
	s.setMethodMetrics = append(s.setMethodMetrics, metric)
	s.setMethodNumbers = append(s.setMethodNumbers, n)
	if labels[0] != nil {
		s.setMethodLabels = append(s.setMethodLabels, labels[0])
	}
}

// Mock StateFactory object
type mockStateFactory struct {
	s *mockState

	createMethodCalled int
}

func newMockStateFactory() *mockStateFactory {
	return &mockStateFactory{createMethodCalled: 0}
}

func (m *mockStateFactory) Create() state.State {
	m.createMethodCalled++
	if m.s == nil {
		m.s = &mockState{
			setMethodCalled:  0,
			setMethodMetrics: []string{},
			setMethodNumbers: []int{},
			setMethodLabels:  []map[string]string{},
		}
	}
	return m.s
}

// Tests block begin

func TestPrometheusClient_NewPrometheus(t *testing.T) {
	p := NewPrometheus("namespace", newMockIncrementerFactory(), newMockStateFactory())
	assert.IsType(t, &Prometheus{}, p)
	assert.Equal(t, "namespace", p.namespace)
	assert.IsType(t, &mockIncrementerFactory{}, p.incFactory)
	assert.IsType(t, &mockStateFactory{}, p.stFactory)
	assert.IsType(t, map[string]incrementer.Incrementer{}, p.increments)
	assert.IsType(t, map[string]state.State{}, p.states)
	assert.IsType(t, map[string]*prometheus.HistogramVec{}, p.histograms)
}

func TestPrometheusClient_BuildTimer(t *testing.T) {
	p := NewPrometheus("namespace", newMockIncrementerFactory(), newMockStateFactory())
	tt := p.BuildTimer()
	_, ok := tt.(*timer.Memory)
	assert.True(t, ok)
}

func TestPrometheusClient_Close(t *testing.T) {
	p := NewPrometheus("namespace", newMockIncrementerFactory(), newMockStateFactory())
	tt := p.Close()
	assert.Nil(t, tt)
}

func TestPrometheusClient_TrackMetric(t *testing.T) {
	m := newMockIncrementerFactory()
	p := NewPrometheus("namespace", m, newMockStateFactory())
	p.TrackMetric("section", bucket.NewMetricOperation("foo", "bar", "baz"))

	assert.Equal(t, 2, m.createMethodCalled)

	assert.Equal(t, 2, m.inc.incrementMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz", "namespace_total_section"}, m.inc.incrementMethodMetrics)
	assert.Equal(t, []map[string]string{}, m.inc.incrementMethodLabels)
}

func TestPrometheusClient_TrackMetricIncrementsAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	p := NewPrometheus("namespace", m, newMockStateFactory())
	p.TrackMetric("section", bucket.NewMetricOperation("foo", "bar", "baz"))
	p.TrackMetric("section", bucket.NewMetricOperation("foo", "bar", "baz"))

	assert.Equal(t, 2, m.createMethodCalled)

	assert.Equal(t, 4, m.inc.incrementMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz", "namespace_total_section", "namespace_section_foo_bar_baz", "namespace_total_section"}, m.inc.incrementMethodMetrics)
	assert.Equal(t, []map[string]string{}, m.inc.incrementMethodLabels)

	assert.Equal(t, 0, m.inc.incrementNMethodCalled)

	assert.Equal(t, 0, m.inc.incrementALlMethodCalled)
	assert.Equal(t, 0, m.inc.incrementALlNMethodCalled)
}

func TestPrometheusClient_TrackMetricN(t *testing.T) {
	m := newMockIncrementerFactory()
	p := NewPrometheus("namespace", m, newMockStateFactory())
	p.TrackMetricN("section", bucket.NewMetricOperation("foo", "bar", "baz"), 999)

	assert.Equal(t, 2, m.createMethodCalled)

	assert.Equal(t, 0, m.inc.incrementMethodCalled)

	assert.Equal(t, 2, m.inc.incrementNMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz", "namespace_total_section"}, m.inc.incrementNMethodMetrics)
	assert.Equal(t, []map[string]string{}, m.inc.incrementNMethodLabels)

	assert.Equal(t, 0, m.inc.incrementALlMethodCalled)
	assert.Equal(t, 0, m.inc.incrementALlNMethodCalled)
}

func TestPrometheusClient_TrackState(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)
	p.TrackState("section", bucket.NewMetricOperation("foo", "bar", "baz"), 888)
	assert.Equal(t, 1, s.createMethodCalled)

	assert.Equal(t, 1, s.s.setMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz"}, s.s.setMethodMetrics)
	assert.Equal(t, []int{888}, s.s.setMethodNumbers)
	assert.Equal(t, []map[string]string{}, s.s.setMethodLabels)
}

func TestPrometheusClient_TrackStateAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)
	p.TrackState("section", bucket.NewMetricOperation("foo", "bar", "baz"), 888)
	p.TrackState("section", bucket.NewMetricOperation("foo", "bar", "baz"), 888)
	assert.Equal(t, 1, s.createMethodCalled)

	assert.Equal(t, 2, s.s.setMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz", "namespace_section_foo_bar_baz"}, s.s.setMethodMetrics)
	assert.Equal(t, []int{888, 888}, s.s.setMethodNumbers)
	assert.Equal(t, []map[string]string{}, s.s.setMethodLabels)
}

func TestPrometheusClient_TrackOperation(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, true)

	assert.Equal(t, 2, m.createMethodCalled)

	assert.Equal(t, 2, m.inc.incrementMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz", "namespace_total_section"}, m.inc.incrementMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}}, m.inc.incrementMethodLabels)
}

func TestPrometheusClient_TrackOperationWithTimer(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()

	p := NewPrometheus("namespace", m, s)

	tt := p.BuildTimer()

	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), tt, true)

	assert.Equal(t, 2, m.createMethodCalled)

	assert.Equal(t, 2, m.inc.incrementMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz", "namespace_total_section"}, m.inc.incrementMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}}, m.inc.incrementMethodLabels)

	assert.NotNil(t, p.histograms["namespace_section_foo_bar_baz"])
	_, err := p.histograms["namespace_section_foo_bar_baz"].GetMetricWithLabelValues("true")
	assert.Nil(t, err)
}

func TestPrometheusClient_TrackOperationAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, true)
	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, false)

	assert.Equal(t, 2, m.createMethodCalled)

	assert.Equal(t, 4, m.inc.incrementMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz", "namespace_total_section", "namespace_section_foo_bar_baz", "namespace_total_section"}, m.inc.incrementMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}, {"success": "false"}, {"success": "false"}}, m.inc.incrementMethodLabels)
}

func TestPrometheusClient_TrackOperationN(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, 999, true)

	assert.Equal(t, 2, m.createMethodCalled)

	assert.Equal(t, 2, m.inc.incrementNMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz", "namespace_total_section"}, m.inc.incrementNMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}}, m.inc.incrementNMethodLabels)
}

func TestPrometheusClient_TrackOperationNWithTimer(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()

	p := NewPrometheus("namespace", m, s)

	tt := p.BuildTimer()

	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), tt, 999, true)

	assert.Equal(t, 2, m.createMethodCalled)
	assert.Equal(t, 2, m.inc.incrementNMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz", "namespace_total_section"}, m.inc.incrementNMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}}, m.inc.incrementNMethodLabels)

	assert.NotNil(t, p.histograms["namespace_section_foo_bar_baz"])
	_, err := p.histograms["namespace_section_foo_bar_baz"].GetMetricWithLabelValues("true")
	assert.Nil(t, err)

}

func TestPrometheusClient_TrackOperationNAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, 1, true)
	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, 2, false)

	assert.Equal(t, 2, m.createMethodCalled)

	assert.Equal(t, 4, m.inc.incrementNMethodCalled)
	assert.Equal(t, []string{"namespace_section_foo_bar_baz", "namespace_total_section", "namespace_section_foo_bar_baz", "namespace_total_section"}, m.inc.incrementNMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}, {"success": "false"}, {"success": "false"}}, m.inc.incrementNMethodLabels)
}
