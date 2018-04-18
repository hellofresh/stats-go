package client

import (
	"testing"
	"time"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/incrementer"
	"github.com/hellofresh/stats-go/state"
	"github.com/hellofresh/stats-go/timer"
	"github.com/stretchr/testify/assert"
)

// Mock incrementer object
type mockIncrementer struct {
	IncrementMethodCalled     int
	IncrementNMethodCalled    int
	IncrementALlMethodCalled  int
	IncrementALlNMethodCalled int
}

func (i *mockIncrementer) Increment(metric string, labels ...map[string]string) {
	i.IncrementMethodCalled++
}

func (i *mockIncrementer) IncrementN(metric string, n int, labels ...map[string]string) {
	i.IncrementNMethodCalled++
}

func (i *mockIncrementer) IncrementAll(b bucket.Bucket) {
	i.IncrementALlMethodCalled++
}

func (i *mockIncrementer) IncrementAllN(b bucket.Bucket, n int) {
	i.IncrementALlNMethodCalled++
}

// Mock Factory object
type mockIncrementerFactory struct {
	M *mockIncrementer

	CreateMethodCalled int
}

func newMockIncrementerFactory() *mockIncrementerFactory {
	return &mockIncrementerFactory{CreateMethodCalled: 0}
}

func (m *mockIncrementerFactory) Create() incrementer.Incrementer {
	m.CreateMethodCalled++
	if m.M == nil {
		m.M = &mockIncrementer{
			IncrementMethodCalled:     0,
			IncrementNMethodCalled:    0,
			IncrementALlMethodCalled:  0,
			IncrementALlNMethodCalled: 0,
		}
	}
	return m.M
}

// Mock state object
type mockState struct {
	SetMethodCalled int
}

func (s *mockState) Set(metric string, n int, labels ...map[string]string) {
	s.SetMethodCalled++
}

// Mock Factory object
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
			SetMethodCalled: 0,
		}
	}
	return m.S
}

// Mock Timer object
type mockTimer struct {
	StartMethodCalled   int
	StartAtMethodCalled int
	FinishMethodCalled  int
}

func newMockTimer() *mockTimer {
	return &mockTimer{
		StartMethodCalled:   0,
		StartAtMethodCalled: 0,
		FinishMethodCalled:  0,
	}
}

func (t *mockTimer) Start() timer.Timer {
	t.StartMethodCalled++
	return t
}

func (t *mockTimer) StartAt(time.Time) timer.Timer {
	t.StartAtMethodCalled++
	return t
}

func (t *mockTimer) Finish(bucket string, labels ...map[string]string) {
	t.FinishMethodCalled++
}

// Tests block begin

func TestPrometheusClient_NewPrometheus(t *testing.T) {
	p := NewPrometheus("namespace", newMockIncrementerFactory(), newMockStateFactory())
	assert.IsType(t, &Prometheus{}, p)
}

func TestPrometheusClient_BuildTimer(t *testing.T) {
	p := NewPrometheus("namespace", newMockIncrementerFactory(), newMockStateFactory())
	tt := p.BuildTimer()
	_, ok := tt.(*timer.Prometheus)
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
	assert.Equal(t, 2, m.M.IncrementMethodCalled)
}

func TestPrometheusClient_TrackMetricIncrementsAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	p := NewPrometheus("namespace", m, newMockStateFactory())
	p.TrackMetric("section", bucket.NewMetricOperation("foo", "bar", "baz"))
	p.TrackMetric("section", bucket.NewMetricOperation("foo", "bar", "baz"))

	assert.Equal(t, 2, m.CreateMethodCalled)
	assert.Equal(t, 4, m.M.IncrementMethodCalled)
	assert.Equal(t, 0, m.M.IncrementNMethodCalled)
	assert.Equal(t, 0, m.M.IncrementALlMethodCalled)
	assert.Equal(t, 0, m.M.IncrementALlNMethodCalled)
}

func TestPrometheusClient_TrackMetricN(t *testing.T) {
	m := newMockIncrementerFactory()
	p := NewPrometheus("namespace", m, newMockStateFactory())
	p.TrackMetricN("section", bucket.NewMetricOperation("foo", "bar", "baz"), 999)

	assert.Equal(t, 2, m.CreateMethodCalled)
	assert.Equal(t, 0, m.M.IncrementMethodCalled)
	assert.Equal(t, 2, m.M.IncrementNMethodCalled)
	assert.Equal(t, 0, m.M.IncrementALlMethodCalled)
	assert.Equal(t, 0, m.M.IncrementALlNMethodCalled)
}

func TestPrometheusClient_TrackState(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)
	p.TrackState("section", bucket.NewMetricOperation("foo", "bar", "baz"), 888)
	assert.Equal(t, 1, s.CreateMethodCalled)
	assert.Equal(t, 1, s.S.SetMethodCalled)
}

func TestPrometheusClient_TrackStateAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)
	p.TrackState("section", bucket.NewMetricOperation("foo", "bar", "baz"), 888)
	p.TrackState("section", bucket.NewMetricOperation("foo", "bar", "baz"), 888)
	assert.Equal(t, 1, s.CreateMethodCalled)
	assert.Equal(t, 2, s.S.SetMethodCalled)
}

func TestPrometheusClient_TrackOperation(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, true)

	assert.Equal(t, 2, m.CreateMethodCalled)
	assert.Equal(t, 2, m.M.IncrementMethodCalled)
}

func TestPrometheusClient_TrackOperationWithTimer(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	mockedTimer := newMockTimer()

	p := NewPrometheus("namespace", m, s)

	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), mockedTimer, true)

	assert.Equal(t, 2, m.CreateMethodCalled)
	assert.Equal(t, 1, mockedTimer.FinishMethodCalled)
	assert.Equal(t, 2, m.M.IncrementMethodCalled)
}

func TestPrometheusClient_TrackOperationAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, true)
	p.TrackOperation("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, false)

	assert.Equal(t, 2, m.CreateMethodCalled)
	assert.Equal(t, 4, m.M.IncrementMethodCalled)
}

func TestPrometheusClient_TrackOperationN(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, 999, true)

	assert.Equal(t, 2, m.CreateMethodCalled)
	assert.Equal(t, 2, m.M.IncrementNMethodCalled)
}

func TestPrometheusClient_TrackOperationNWithTimer(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	mockedTimer := newMockTimer()

	p := NewPrometheus("namespace", m, s)

	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), mockedTimer, 999, true)

	assert.Equal(t, 2, m.CreateMethodCalled)
	assert.Equal(t, 1, mockedTimer.FinishMethodCalled)
	assert.Equal(t, 2, m.M.IncrementNMethodCalled)
}

func TestPrometheusClient_TrackOperationNAlreadyExists(t *testing.T) {
	m := newMockIncrementerFactory()
	s := newMockStateFactory()
	p := NewPrometheus("namespace", m, s)

	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, 1, true)
	p.TrackOperationN("section", bucket.NewMetricOperation("foo", "bar", "baz"), nil, 2, false)

	assert.Equal(t, 2, m.CreateMethodCalled)
	assert.Equal(t, 4, m.M.IncrementNMethodCalled)
}
