package client

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/stretchr/testify/assert"
)

func TestNewLogClient_SetHTTPMetricCallback(t *testing.T) {
	client := NewLog(true)
	callback := func(metricParts *bucket.MetricOperation, r *http.Request) *bucket.MetricOperation {
		return metricParts
	}

	client.SetHTTPMetricCallback(callback)
	// asserting functions directly gives false result:
	// Not equal: (func(bucket.MetricOperation, *http.Request) bucket.MetricOperation)(0x1255160) (expected)
	//     != (bucket.HTTPMetricNameAlterCallback)(0x1255160) (actual)
	// so we assert objects pointers directly to make sure this is the same function object
	assert.Equal(t, reflect.ValueOf(callback).Pointer(), reflect.ValueOf(client.GetHTTPMetricCallback()).Pointer())
}
