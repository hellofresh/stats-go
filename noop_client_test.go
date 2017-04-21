package stats

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopClient(t *testing.T) {
	client := NewNoopClient()

	assert.Nil(t, client.Close())
	assert.IsType(t, &MemoryTimeTracker{}, client.BuildTimeTracker())
	assert.Equal(t, client, client.TrackRequest(nil, nil, true))
	assert.Equal(t, client, client.TrackOperation("", MetricOperation{}, nil, true))
	assert.Equal(t, client, client.TrackOperationN("", MetricOperation{}, nil, 0, true))
	assert.Equal(t, client, client.SetHTTPRequestSection(""))
	assert.Equal(t, client, client.ResetHTTPRequestSection())
	assert.Equal(t, client, client.SetHTTPMetricCallback(func(metricParts MetricOperation, r *http.Request) MetricOperation {
		return MetricOperation{}
	}))
}
