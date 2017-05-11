package stats

import (
	"net/http"
	"testing"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/timer"
	"github.com/stretchr/testify/assert"
)

func TestNoopClient(t *testing.T) {
	client := NewNoopClient()

	assert.Nil(t, client.Close())
	assert.IsType(t, &timer.Memory{}, client.BuildTimer())
	assert.Equal(t, client, client.TrackRequest(nil, nil, true))
	assert.Equal(t, client, client.TrackOperation("", bucket.MetricOperation{}, nil, true))
	assert.Equal(t, client, client.TrackOperationN("", bucket.MetricOperation{}, nil, 0, true))
	assert.Equal(t, client, client.SetHTTPRequestSection(""))
	assert.Equal(t, client, client.ResetHTTPRequestSection())
	assert.Equal(t, client, client.SetHTTPMetricCallback(func(metricParts bucket.MetricOperation, r *http.Request) bucket.MetricOperation {
		return bucket.MetricOperation{}
	}))
}
