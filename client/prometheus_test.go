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
	IncrementMethodCalled  int
	IncrementMethodMetrics []string
	IncrementMethodLabels  []map[string]string

	IncrementNMethodCalled  int
	IncrementNMethodMetrics []string
	IncrementNMethodLabels  []map[string]string

	IncrementALlMethodCalled  int
	IncrementALlNMethodCalled int
}

func (i *mockIncrementer) Increment(metric string, labels ...map[string]string) {
	i.IncrementMethodCalled++
	i.IncrementMethodMetrics = append(i.IncrementMethodMetrics, metric)
	if labels[0] != nil {
		i.IncrementMethodLabels = append(i.IncrementMethodLabels, labels[0])
	}
}

func (i *mockIncrementer) IncrementN(metric string, n int, labels ...map[string]string) {
	i.IncrementNMethodCalled++
	i.IncrementNMethodMetrics = append(i.IncrementNMethodMetrics, metric)
	if labels[0] != nil {
		i.IncrementNMethodLabels = append(i.IncrementNMethodLabels, labels[0])
	}
}

func (i *mockIncrementer) IncrementAll(b bucket.Bucket) {
	i.IncrementALlMethodCalled++
}

func (i *mockIncrementer) IncrementAllN(b bucket.Bucket, n int) {
	i.IncrementALlNMethodCalled++
}

// Mock IncrementFactory object
type mockIncrementerFactory struct {
	inc *mockIncrementer

	CreateMethodCalled int
}

func newMockIncrementerFactory() *mockIncrementerFactory {
	return &mockIncrementerFactory{CreateMethodCalled: 0}
}

func (m *mockIncrementerFactory) Create() incrementer.Incrementer {
	m.CreateMethodCalled++
	if m.inc == nil {
		m.inc = &mockIncrementer{
			IncrementMethodCalled:  0,
			IncrementMethodMetrics: []string{},
			IncrementMethodLabels:  []map[string]string{},

			IncrementNMethodCalled:  0,
			IncrementNMethodMetrics: []string{},
			IncrementNMethodLabels:  []map[string]string{},

			IncrementALlMethodCalled:  0,
			IncrementALlNMethodCalled: 0,
		}
	}
	return m.inc
}

// Mock state object
type mockState struct {
	SetMethodCalled  int
	SetMethodNumbers []int
	SetMethodMetrics []string
	SetMethodLabels  []map[string]string
}

func (s *mockState) Set(metric string, n int, labels ...map[string]string) {
	s.SetMethodCalled++
	s.SetMethodMetrics = append(s.SetMethodMetrics, metric)
	s.SetMethodNumbers = append(s.SetMethodNumbers, n)
	if labels[0] != nil {
		s.SetMethodLabels = append(s.SetMethodLabels, labels[0])
	}
}

// Mock StateFactory object
type mockStateFactory struct {
	S *mockState

	CreateMethodCalled int
}

func newMockStateFactory() *mockStateFactory {
	return &mockStateFactory{CreateMethodCalled: 0}
}

func (m *mockStateFactory) Create() state.State {
	m.CreateMethodCalled++
	if m.S == nil {
		m.S = &mockState{
			SetMethodCalled:  0,
			SetMethodMetrics: []string{},
			SetMethodNumbers: []int{},
			SetMethodLabels:  []map[string]string{},
		}
	}
	return m.S
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

	assert.Equal(t, 2, m.CreateMethodCalled)

	assert.Equal(t, 2, m.inc.IncrementMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz", "total_section"}, m.inc.IncrementMethodMetrics)
	assert.Equal(t, []map[string]string{}, m.inc.IncrementMethodLabels)
}

func TestPrometheusClient_TrackMetricIncrementsAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	p := NewPrometheus("namespace", m, newMockStateFactory())
	p.TrackMetric("section", bucket.NewMetricOperation("foo", "bar", "baz"))
	p.TrackMetric("section", bucket.NewMetricOperation("foo", "bar", "baz"))

	assert.Equal(t, 2, m.CreateMethodCalled)

	assert.Equal(t, 4, m.inc.IncrementMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz", "total_section", "section_foo_bar_baz", "total_section"}, m.inc.IncrementMethodMetrics)
	assert.Equal(t, []map[string]string{}, m.inc.IncrementMethodLabels)

	assert.Equal(t, 0, m.inc.IncrementNMethodCalled)

	assert.Equal(t, 0, m.inc.IncrementALlMethodCalled)
	assert.Equal(t, 0, m.inc.IncrementALlNMethodCalled)
}

func TestPrometheusClient_TrackMetricN(t *testing.T) {
	m := newMockIncrementerFactory()
	p := NewPrometheus("namespace", m, newMockStateFactory())
	p.TrackMetricN("section", bucket.NewMetricOperation("foo", "bar", "baz"), 999)

	assert.Equal(t, 2, m.CreateMethodCalled)

	assert.Equal(t, 0, m.inc.IncrementMethodCalled)

	assert.Equal(t, 2, m.inc.IncrementNMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz", "total_section"}, m.inc.IncrementNMethodMetrics)
	assert.Equal(t, []map[string]string{}, m.inc.IncrementNMethodLabels)

	assert.Equal(t, 0, m.inc.IncrementALlMethodCalled)
	assert.Equal(t, 0, m.inc.IncrementALlNMethodCalled)
}

func TestPrometheusClient_TrackState(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)
	p.TrackState("section", bucket.NewMetricOperation("foo", "bar", "baz"), 888)
	assert.Equal(t, 1, s.CreateMethodCalled)

	assert.Equal(t, 1, s.S.SetMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz"}, s.S.SetMethodMetrics)
	assert.Equal(t, []int{888}, s.S.SetMethodNumbers)
	assert.Equal(t, []map[string]string{}, s.S.SetMethodLabels)
}

func TestPrometheusClient_TrackStateAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)
	p.TrackState("section", bucket.NewMetricOperation("foo", "bar", "baz"), 888)
	p.TrackState("section", bucket.NewMetricOperation("foo", "bar", "baz"), 888)
	assert.Equal(t, 1, s.CreateMethodCalled)

	assert.Equal(t, 2, s.S.SetMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz", "section_foo_bar_baz"}, s.S.SetMethodMetrics)
	assert.Equal(t, []int{888, 888}, s.S.SetMethodNumbers)
	assert.Equal(t, []map[string]string{}, s.S.SetMethodLabels)
}

func TestPrometheusClient_TrackOperation(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, true)

	assert.Equal(t, 2, m.CreateMethodCalled)

	assert.Equal(t, 2, m.inc.IncrementMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz", "total_section"}, m.inc.IncrementMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}}, m.inc.IncrementMethodLabels)
}

func TestPrometheusClient_TrackOperationWithTimer(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()

	p := NewPrometheus("namespace", m, s)

	tt := p.BuildTimer()

	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), tt, true)

	assert.Equal(t, 2, m.CreateMethodCalled)

	assert.Equal(t, 2, m.inc.IncrementMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz", "total_section"}, m.inc.IncrementMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}}, m.inc.IncrementMethodLabels)

	assert.NotNil(t, p.histograms["section_foo_bar_baz"])
	_, err := p.histograms["section_foo_bar_baz"].GetMetricWithLabelValues("true")
	assert.Nil(t, err)
}

func TestPrometheusClient_TrackOperationAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, true)
	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, false)

	assert.Equal(t, 2, m.CreateMethodCalled)

	assert.Equal(t, 4, m.inc.IncrementMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz", "total_section", "section_foo_bar_baz", "total_section"}, m.inc.IncrementMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}, {"success": "false"}, {"success": "false"}}, m.inc.IncrementMethodLabels)
}

func TestPrometheusClient_TrackOperationN(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, 999, true)

	assert.Equal(t, 2, m.CreateMethodCalled)

	assert.Equal(t, 2, m.inc.IncrementNMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz", "total_section"}, m.inc.IncrementNMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}}, m.inc.IncrementNMethodLabels)
}

func TestPrometheusClient_TrackOperationNWithTimer(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()

	p := NewPrometheus("namespace", m, s)

	tt := p.BuildTimer()

	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), tt, 999, true)

	assert.Equal(t, 2, m.CreateMethodCalled)
	assert.Equal(t, 2, m.inc.IncrementNMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz", "total_section"}, m.inc.IncrementNMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}}, m.inc.IncrementNMethodLabels)

	assert.NotNil(t, p.histograms["section_foo_bar_baz"])
	_, err := p.histograms["section_foo_bar_baz"].GetMetricWithLabelValues("true")
	assert.Nil(t, err)

}

func TestPrometheusClient_TrackOperationNAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, 1, true)
	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, 2, false)

	assert.Equal(t, 2, m.CreateMethodCalled)

	assert.Equal(t, 4, m.inc.IncrementNMethodCalled)
	assert.Equal(t, []string{"section_foo_bar_baz", "total_section", "section_foo_bar_baz", "total_section"}, m.inc.IncrementNMethodMetrics)
	assert.Equal(t, []map[string]string{{"success": "true"}, {"success": "true"}, {"success": "false"}, {"success": "false"}}, m.inc.IncrementNMethodLabels)
}
