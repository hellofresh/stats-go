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
	incFactory incrementer.Factory
	stFactory  state.Factory

	increments map[string]incrementer.Incrementer
	states     map[string]state.State
}

// NewPrometheus builds and returns new Prometheus instance
func NewPrometheus(namespace string, incFactory incrementer.Factory, stFactory state.Factory) *Prometheus {
	client := &Prometheus{
		namespace:  namespace,
		incFactory: incFactory,
		stFactory:  stFactory,
		increments: make(map[string]incrementer.Incrementer),
		states:     make(map[string]state.State),
	}
	return client
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
	b := bucket.NewHTTPRequest(c.httpRequestSection, r, success, c.httpMetricCallback, c.unicode)
	metric := b.Metric()
	metricTotal := b.MetricTotal()

	metric = strings.Replace(metric, "-.", "", -1)
	metric = strings.Replace(metric, ".-", "", -1)
	metric = strings.Replace(metric, "-", "", -1)
	metric = strings.Replace(metric, ".", "_", -1)

	metricTotal = strings.Replace(metricTotal, "-.", "", -1)
	metricTotal = strings.Replace(metricTotal, ".-", "", -1)
	metricTotal = strings.Replace(metricTotal, "-", "", -1)
	metricTotal = strings.Replace(metricTotal, ".", "_", -1)

	c.Lock()
	defer c.Unlock()

	if _, ok := c.increments[metric]; !ok {
		c.increments[metric] = c.incFactory.Create()
	}

	if _, ok := c.increments[metricTotal]; !ok {
		c.increments[metricTotal] = c.incFactory.Create()
	}

	labels := map[string]string{"success": strconv.FormatBool(success), "action": r.Method}

	c.increments[metric].Increment(metric, labels)
	c.increments[metricTotal].Increment(metricTotal, labels)

	return c
}

// TrackOperation tracks custom operation
func (c *Prometheus) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	b := bucket.NewPrometheus(section, operation, success, c.unicode)

	if operation.Labels == nil {
		operation.Labels = map[string]string{"success": strconv.FormatBool(success)}
	} else {
		operation.Labels["success"] = strconv.FormatBool(success)
	}

	c.TrackMetric(section, operation)

	if nil != t {
		t.Finish(b.Metric(), operation.Labels)
	}

	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *Prometheus) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	b := bucket.NewPrometheus(section, operation, success, c.unicode)

	if operation.Labels == nil {
		operation.Labels = map[string]string{"success": strconv.FormatBool(success)}
	} else {
		operation.Labels["success"] = strconv.FormatBool(success)
	}

	c.TrackMetricN(section, operation, n)

	if nil != t {
		t.Finish(b.Metric(), operation.Labels)
	}
	return c
}

// TrackMetric tracks custom metric, w/out ok/fail additional sections
func (c *Prometheus) TrackMetric(section string, operation bucket.MetricOperation) Client {
	b := bucket.NewPrometheus(section, operation, true, c.unicode)
	metric := b.Metric()
	metricTotal := b.MetricTotal()

	c.Lock()
	defer c.Unlock()

	if _, ok := c.increments[metric]; !ok {
		c.increments[metric] = c.incFactory.Create()
	}

	if _, ok := c.increments[metricTotal]; !ok {
		c.increments[metricTotal] = c.incFactory.Create()
	}
	c.increments[metric].Increment(metric, operation.Labels)
	c.increments[metricTotal].Increment(metricTotal, operation.Labels)

	return c
}

// TrackMetricN tracks custom metric with n diff, w/out ok/fail additional sections
func (c *Prometheus) TrackMetricN(section string, operation bucket.MetricOperation, n int) Client {
	b := bucket.NewPrometheus(section, operation, true, c.unicode)
	metric := b.Metric()
	metricTotal := b.MetricTotal()

	c.Lock()
	defer c.Unlock()

	if _, ok := c.increments[metric]; !ok {
		c.increments[metric] = c.incFactory.Create()
	}

	if _, ok := c.increments[metricTotal]; !ok {
		c.increments[metricTotal] = c.incFactory.Create()
	}
	c.increments[metric].IncrementN(metric, n, operation.Labels)
	c.increments[metricTotal].IncrementN(metricTotal, n, operation.Labels)

	return c
}

// TrackState tracks metric absolute value
func (c *Prometheus) TrackState(section string, operation bucket.MetricOperation, value int) Client {
	b := bucket.NewPrometheus(section, operation, true, c.unicode)
	metric := b.Metric()

	c.Lock()
	defer c.Unlock()

	if _, ok := c.states[metric]; !ok {
		c.states[metric] = c.stFactory.Create()
	}

	c.states[metric].Set(metric, value, operation.Labels)

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
