package stats

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryClient_BuildTimeTracker(t *testing.T) {
	client := NewMemoryClient()
	tt := client.BuildTimeTracker()
	_, ok := tt.(*MemoryTimeTracker)
	assert.True(t, ok)
}

func TestMemoryClient_TrackRequest(t *testing.T) {
	client := NewMemoryClient()

	tt := client.BuildTimeTracker()
	r := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/hello/memory/test"}}
	success := true
	b := NewBucketHTTPRequest(client.httpRequestSection, r, success, client.httpMetricCallback)

	client.TrackRequest(r, tt, success)

	assert.Equal(t, 1, len(client.TimeMetrics))
	assert.Equal(t, 4, len(client.CountMetrics))

	assert.Equal(t, b.Metric(), client.TimeMetrics[0].Bucket)
	assert.Equal(t, 1, client.CountMetrics[b.Metric()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricWithSuffix()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricTotal()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricTotalWithSuffix()])

	client.Close()

	assert.Equal(t, 0, len(client.TimeMetrics))
	assert.Equal(t, 0, len(client.CountMetrics))
}

func TestMemoryClient_TrackOperation(t *testing.T) {
	client := NewMemoryClient()

	tt := client.BuildTimeTracker()
	section := "test-section"
	operation := MetricOperation{"o1", "o2", "o3"}
	success := true
	b := NewBucketPlain(section, operation, success)

	client.TrackOperation(section, operation, tt, success)

	assert.Equal(t, 1, len(client.TimeMetrics))
	assert.Equal(t, 4, len(client.CountMetrics))

	assert.Equal(t, b.MetricWithSuffix(), client.TimeMetrics[0].Bucket)
	assert.Equal(t, 1, client.CountMetrics[b.Metric()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricWithSuffix()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricTotal()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricTotalWithSuffix()])

	client.Close()

	assert.Equal(t, 0, len(client.TimeMetrics))
	assert.Equal(t, 0, len(client.CountMetrics))
}

func TestMemoryClient_TrackOperationN(t *testing.T) {
	client := NewMemoryClient()

	tt := client.BuildTimeTracker()
	section := "test-section"
	operation := MetricOperation{"o1", "o2", "o3"}
	success := true
	n := 5
	b := NewBucketPlain(section, operation, success)

	client.TrackOperationN(section, operation, tt, n, success)

	assert.Equal(t, 1, len(client.TimeMetrics))
	assert.Equal(t, 4, len(client.CountMetrics))

	assert.Equal(t, b.MetricWithSuffix(), client.TimeMetrics[0].Bucket)
	assert.Equal(t, n, client.CountMetrics[b.Metric()])
	assert.Equal(t, n, client.CountMetrics[b.MetricWithSuffix()])
	assert.Equal(t, n, client.CountMetrics[b.MetricTotal()])
	assert.Equal(t, n, client.CountMetrics[b.MetricTotalWithSuffix()])

	client.Close()

	assert.Equal(t, 0, len(client.TimeMetrics))
	assert.Equal(t, 0, len(client.CountMetrics))
}

func TestMemoryClient_SetHTTPRequestSection(t *testing.T) {
	client := NewMemoryClient()

	assert.Equal(t, sectionRequest, client.httpRequestSection)

	section := "test-section"
	client.SetHTTPRequestSection(section)
	assert.Equal(t, section, client.httpRequestSection)

	client.ResetHTTPRequestSection()
	assert.Equal(t, sectionRequest, client.httpRequestSection)
}
