package client

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/incrementer"
	"github.com/hellofresh/stats-go/state"
	"github.com/hellofresh/stats-go/timer"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus is Client implementation for prometheus
type Prometheus struct {
	sync.Mutex

	unicode            bool
	httpMetricCallback bucket.HTTPMetricNameAlterCallback
	httpRequestSection string

	namespace  string
	increments map[string]*incrementer.Prometheus
	states     map[string]*state.Prometheus
}

// NewPrometheus builds and returns new Prometheus instance
func NewPrometheus(namespace string) (*Prometheus, error) {
	client := &Prometheus{
		namespace:  namespace,
		increments: make(map[string]*incrementer.Prometheus),
		states:     make(map[string]*state.Prometheus),
	}
	return client, nil
}

// BuildTimer builds timer to track metric timings
func (c *Prometheus) BuildTimer() timer.Timer {
	return timer.NewPrometheus()
}

// Close closes underlying client connection if any
func (c *Prometheus) Close() error {
	return nil
}

// TrackRequest tracks HTTP Request stats
func (c *Prometheus) TrackRequest(r *http.Request, t timer.Timer, success bool) Client {
	c.Lock()
	defer c.Unlock()

	b := bucket.NewHTTPRequest(c.httpRequestSection, r, success, c.httpMetricCallback, c.unicode)
	metric := b.Metric()
	metricTotal := b.MetricTotal()

	// need to do smth with that
	metric = strings.Replace(metric, "-.", "", -1)
	metric = strings.Replace(metric, ".-", "", -1)
	metric = strings.Replace(metric, "-", "", -1)
	metric = strings.Replace(metric, ".", "_", -1)

	metricTotal = strings.Replace(metricTotal, "-.", "", -1)
	metricTotal = strings.Replace(metricTotal, ".-", "", -1)
	metricTotal = strings.Replace(metricTotal, "-", "", -1)
	metricTotal = strings.Replace(metricTotal, ".", "_", -1)

	if _, ok := c.increments[metric]; !ok {
		c.increments[metric] = incrementer.NewPrometheus(incrementer.NewPrometheusCounterFactory())
	}

	if _, ok := c.increments[metricTotal]; !ok {
		c.increments[metricTotal] = incrementer.NewPrometheus(incrementer.NewPrometheusCounterFactory())
	}

	labelNames := []string{"success", "action"}
	labelValues := []string{strconv.FormatBool(success), r.Method}

	c.increments[metric].IncrementWithLabels(metric, labelNames, labelValues)
	c.increments[metricTotal].IncrementWithLabels(metricTotal, labelNames, labelValues)

	return c
}

// TrackOperation tracks custom operation
func (c *Prometheus) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	b := bucket.NewPrometheus(section, operation, success, c.unicode)
	if nil != t {
		t.FinishWithLabels(b.Metric(), map[string]string{"success": strconv.FormatBool(success)})
	}

	labelValues := []string{strconv.FormatBool(success)}
	for _, v := range operation.LabelValues {
		labelValues = append(labelValues, v)
	}

	c.TrackMetric(section, operation.WithLabels(labelValues...))
	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *Prometheus) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	b := bucket.NewPrometheus(section, operation, success, c.unicode)
	if nil != t {
		t.FinishWithLabels(b.Metric(), map[string]string{"success": strconv.FormatBool(success)})
	}
	labelValues := []string{strconv.FormatBool(success)}
	for _, v := range operation.LabelValues {
		labelValues = append(labelValues, v)
	}

	c.TrackMetricN(section, operation.WithLabels(labelValues...), n)
	return c
}

// TrackMetric tracks custom metric, w/out ok/fail additional sections
func (c *Prometheus) TrackMetric(section string, operation bucket.MetricOperation) Client {
	c.Lock()
	defer c.Unlock()

	b := bucket.NewPrometheus(section, operation, true, c.unicode)
	metric := b.Metric()
	metricTotal := b.MetricTotal()

	if _, ok := c.increments[metric]; !ok {
		c.increments[metric] = incrementer.NewPrometheus(incrementer.NewPrometheusCounterFactory())
	}

	if _, ok := c.increments[metricTotal]; !ok {
		c.increments[metricTotal] = incrementer.NewPrometheus(incrementer.NewPrometheusCounterFactory())
	}

	var labelValues []string

	if len(operation.LabelValues) != len(operation.LabelNames) {
		labelValues = []string{"true"}
		for _, v := range operation.LabelValues {
			labelValues = append(labelValues, v)
		}
	} else {
		labelValues = operation.LabelValues
	}

	c.increments[metric].IncrementWithLabels(metric, operation.LabelNames, labelValues)
	c.increments[metricTotal].IncrementWithLabels(metricTotal, operation.LabelNames, labelValues)

	return c
}

// TrackMetricN tracks custom metric with n diff, w/out ok/fail additional sections
func (c *Prometheus) TrackMetricN(section string, operation bucket.MetricOperation, n int) Client {
	c.Lock()
	defer c.Unlock()

	b := bucket.NewPrometheus(section, operation, true, c.unicode)
	metric := b.Metric()
	metricTotal := b.MetricTotal()

	if _, ok := c.increments[metric]; !ok {
		c.increments[metric] = incrementer.NewPrometheus(incrementer.NewPrometheusCounterFactory())
	}

	if _, ok := c.increments[metricTotal]; !ok {
		c.increments[metricTotal] = incrementer.NewPrometheus(incrementer.NewPrometheusCounterFactory())
	}

	var labelValues []string

	if len(operation.LabelValues) != len(operation.LabelNames) {
		labelValues = []string{"true"}
		for _, v := range operation.LabelValues {
			labelValues = append(labelValues, v)
		}
	} else {
		labelValues = operation.LabelValues
	}

	c.increments[metric].IncrementNWithLabels(metric, n, operation.LabelNames, labelValues)
	c.increments[metricTotal].IncrementNWithLabels(metricTotal, n, operation.LabelNames, labelValues)

	return c
}

// TrackState tracks metric absolute value
func (c *Prometheus) TrackState(section string, operation bucket.MetricOperation, value int) Client {
	c.Lock()
	defer c.Unlock()

	b := bucket.NewPrometheus(section, operation, true, c.unicode)
	metric := b.Metric()

	if _, ok := c.increments[metric]; !ok {
		c.states[metric] = state.NewPrometheus(state.NewPrometheusGaugeFactory())
	}

	labelValues := []string{"true"}
	for _, v := range operation.LabelValues {
		labelValues = append(labelValues, v)
	}

	c.states[metric].SetWithLabels(metric, value, operation.LabelNames, labelValues)

	return c
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (c *Prometheus) SetHTTPMetricCallback(callback bucket.HTTPMetricNameAlterCallback) Client {
	c.Lock()
	defer c.Unlock()

	c.httpMetricCallback = callback
	return c
}

// GetHTTPMetricCallback gets callback handler that allows metric operation alteration for HTTP Request
func (c *Prometheus) GetHTTPMetricCallback() bucket.HTTPMetricNameAlterCallback {
	c.Lock()
	defer c.Unlock()

	return c.httpMetricCallback
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (c *Prometheus) SetHTTPRequestSection(section string) Client {
	c.Lock()
	defer c.Unlock()

	c.httpRequestSection = section
	return c
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (c *Prometheus) ResetHTTPRequestSection() Client {
	return c.SetHTTPRequestSection(bucket.SectionRequest)
}

// Handler returns metrics endpoint for prometheus backend
func (c *Prometheus) Handler() http.Handler {
	return prometheus.Handler()
}
