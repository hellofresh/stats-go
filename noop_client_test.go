package stats

import (
	"net/http"
	"reflect"
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
	assert.Equal(t, client, client.TrackMetric("", bucket.MetricOperation{}))
	assert.Equal(t, client, client.TrackMetricN("", bucket.MetricOperation{}, 0))
	assert.Equal(t, client, client.TrackState("", bucket.MetricOperation{}, 0))
	assert.Equal(t, client, client.SetHTTPRequestSection(""))
	assert.Equal(t, client, client.ResetHTTPRequestSection())
	assert.Equal(t, client, client.SetHTTPMetricCallback(func(metricParts bucket.MetricOperation, r *http.Request) bucket.MetricOperation {
		return bucket.MetricOperation{}
	}))
}

func TestNewNoopClient_SetHTTPMetricCallback(t *testing.T) {
	client := NewNoopClient()
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
