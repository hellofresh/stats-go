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
	labels     map[string]string
}

// NewPrometheus builds and returns new Prometheus instance
func NewPrometheus(namespace string) (*Prometheus, error) {
	client := &Prometheus{
		namespace:  namespace,
		increments: make(map[string]*incrementer.Prometheus),
		states:     make(map[string]*state.Prometheus),
		labels:     make(map[string]string),
	}
	return client, nil
}

// WithLabels adds labels map to metrics
func (c *Prometheus) WithLabels(labels map[string]string) Client {
	c.Lock()
	defer c.Unlock()

	c.labels = labels
	return c
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

	keyHash := ""

	for key := range c.labels {
		keyHash += "_" + key
	}

	if _, ok := c.increments[metric+keyHash]; !ok {
		c.increments[metric+keyHash] = incrementer.NewPrometheus(incrementer.NewPrometheusCounterFactory())
	}

	if _, ok := c.increments[metricTotal+keyHash]; !ok {
		c.increments[metricTotal+keyHash] = incrementer.NewPrometheus(incrementer.NewPrometheusCounterFactory())
	}

	requestLabels := map[string]string{
		"action":  r.Method,
		"success": strconv.FormatBool(success),
	}

	c.increments[metric+keyHash].IncrementWithLabels(metric, requestLabels)
	c.increments[metricTotal+keyHash].IncrementWithLabels(metricTotal, requestLabels)
	c.labels = map[string]string{}

	return c
}

// TrackOperation tracks custom operation
func (c *Prometheus) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	b := bucket.NewPrometheus(section, operation, success, c.unicode)
	if nil != t {
		t.FinishWithLabels(b.Metric(), map[string]string{"success": strconv.FormatBool(success)})
	}
	c.WithLabels(map[string]string{"success": strconv.FormatBool(success)}).TrackMetric(section, operation)
	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *Prometheus) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	b := bucket.NewPrometheus(section, operation, success, c.unicode)
	if nil != t {
		t.FinishWithLabels(b.Metric(), map[string]string{"success": strconv.FormatBool(success)})
	}
	c.WithLabels(map[string]string{"success": strconv.FormatBool(success)}).TrackMetricN(section, operation, n)
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

	c.increments[metric].IncrementWithLabels(metric, c.labels)
	c.increments[metricTotal].IncrementWithLabels(metricTotal, c.labels)
	c.labels = map[string]string{}

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

	c.increments[metric].IncrementNWithLabels(metric, n, c.labels)
	c.increments[metricTotal].IncrementNWithLabels(metricTotal, n, c.labels)
	c.labels = map[string]string{}

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

	c.states[metric].SetWithLabels(metric, value, c.labels)
	c.labels = map[string]string{}

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
