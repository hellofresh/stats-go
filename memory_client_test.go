package stats

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/timer"
	"github.com/stretchr/testify/assert"
)

func TestMemoryClient_BuildTimeTracker(t *testing.T) {
	client := NewMemoryClient()
	tt := client.BuildTimer()
	_, ok := tt.(*timer.Memory)
	assert.True(t, ok)
}

func TestMemoryClient_TrackRequest(t *testing.T) {
	client := NewMemoryClient()

	tt := client.BuildTimer()
	r := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/hello/memory/test"}}
	success := true
	b := bucket.NewHTTPRequest(client.httpRequestSection, r, success, client.httpMetricCallback)

	client.TrackRequest(r, tt, success)

	assert.Equal(t, 1, len(client.TimerMetrics))
	assert.Equal(t, 4, len(client.CountMetrics))

	assert.Equal(t, b.Metric(), client.TimerMetrics[0].Bucket)
	assert.Equal(t, 1, client.CountMetrics[b.Metric()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricWithSuffix()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricTotal()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricTotalWithSuffix()])

	client.Close()

	assert.Equal(t, 0, len(client.TimerMetrics))
	assert.Equal(t, 0, len(client.CountMetrics))
}

func TestMemoryClient_TrackOperation(t *testing.T) {
	client := NewMemoryClient()

	tt := client.BuildTimer()
	section := "test-section"
	operation := bucket.MetricOperation{"o1", "o2", "o3"}
	success := true
	b := bucket.NewPlain(section, operation, success)

	client.TrackOperation(section, operation, tt, success)

	assert.Equal(t, 1, len(client.TimerMetrics))
	assert.Equal(t, 4, len(client.CountMetrics))

	assert.Equal(t, b.MetricWithSuffix(), client.TimerMetrics[0].Bucket)
	assert.Equal(t, 1, client.CountMetrics[b.Metric()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricWithSuffix()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricTotal()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricTotalWithSuffix()])

	client.Close()

	assert.Equal(t, 0, len(client.TimerMetrics))
	assert.Equal(t, 0, len(client.CountMetrics))
}

func TestMemoryClient_TrackOperationN(t *testing.T) {
	client := NewMemoryClient()

	tt := client.BuildTimer()
	section := "test-section"
	operation := bucket.MetricOperation{"o1", "o2", "o3"}
	success := true
	n := 5
	b := bucket.NewPlain(section, operation, success)

	client.TrackOperationN(section, operation, tt, n, success)

	assert.Equal(t, 1, len(client.TimerMetrics))
	assert.Equal(t, 4, len(client.CountMetrics))

	assert.Equal(t, b.MetricWithSuffix(), client.TimerMetrics[0].Bucket)
	assert.Equal(t, n, client.CountMetrics[b.Metric()])
	assert.Equal(t, n, client.CountMetrics[b.MetricWithSuffix()])
	assert.Equal(t, n, client.CountMetrics[b.MetricTotal()])
	assert.Equal(t, n, client.CountMetrics[b.MetricTotalWithSuffix()])

	client.Close()

	assert.Equal(t, 0, len(client.TimerMetrics))
	assert.Equal(t, 0, len(client.CountMetrics))
}

func TestMemoryClient_TrackState(t *testing.T) {
	client := NewMemoryClient()

	section := "test-section"
	operation1 := bucket.MetricOperation{"o1", "o2", "o3"}
	operation2 := bucket.MetricOperation{"p1", "p2", "p3"}
	state1 := 13
	state2 := 66
	state12 := 77

	client.TrackState(section, operation1, state1)
	client.TrackState(section, operation2, state2)

	assert.Equal(t, 2, len(client.StateMetrics))
	assert.Equal(t, state1, client.StateMetrics[bucket.NewPlain(section, operation1, true).Metric()])
	assert.Equal(t, state2, client.StateMetrics[bucket.NewPlain(section, operation2, true).Metric()])

	client.TrackState(section, operation1, state12)
	assert.Equal(t, 2, len(client.StateMetrics))
	assert.Equal(t, state12, client.StateMetrics[bucket.NewPlain(section, operation1, true).Metric()])
	assert.Equal(t, state2, client.StateMetrics[bucket.NewPlain(section, operation2, true).Metric()])
}

func TestMemoryClient_SetHTTPMetricCallback(t *testing.T) {
	client := NewMemoryClient()
	callback := func(metricParts bucket.MetricOperation, r *http.Request) bucket.MetricOperation {
		return metricParts
	}

	client.SetHTTPMetricCallback(callback)
	// asserting functions directly gives false result:
	// Not equal: (func(bucket.MetricOperation, *http.Request) bucket.MetricOperation)(0x1255160) (expected)
	//     != (bucket.HTTPMetricNameAlterCallback)(0x1255160) (actual)
	// so we assert objects pointers directly to make sure this is the same function object
	assert.Equal(t, reflect.ValueOf(callback).Pointer(), reflect.ValueOf(client.GetHTTPMetricCallback()).Pointer())
}

func TestMemoryClient_SetHTTPRequestSection(t *testing.T) {
	client := NewMemoryClient()

	assert.Equal(t, bucket.SectionRequest, client.httpRequestSection)

	section := "test-section"
	client.SetHTTPRequestSection(section)
	assert.Equal(t, section, client.httpRequestSection)

	client.ResetHTTPRequestSection()
	assert.Equal(t, bucket.SectionRequest, client.httpRequestSection)
}
