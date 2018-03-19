package incrementer

import (
	"github.com/hellofresh/stats-go/bucket"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus struct is Incrementer interface implementation that writes all metrics to Prometheus
type Prometheus struct {
	counter        CounterVec
	counterFactory CounterFactory
}

// CounterVec interface for counter vectors in prometheus backend
type CounterVec interface {
	GetMetricWithLabelValues(lvs ...string) (prometheus.Counter, error)
	GetMetricWith(labels prometheus.Labels) (prometheus.Counter, error)
	WithLabelValues(lvs ...string) prometheus.Counter
	With(labels prometheus.Labels) prometheus.Counter
}

// CounterFactory interface for making new CounterVec instances
type CounterFactory interface {
	Create(metric string, labelKeys []string) CounterVec
}

// PrometheusCounterFactory implements CounterFactory interface
type PrometheusCounterFactory struct {
}

// NewPrometheusCounterFactory returns new PrometheusCounterFactory instance
func NewPrometheusCounterFactory() *PrometheusCounterFactory {
	return &PrometheusCounterFactory{}
}

// Create method returns new CounterVec instance with metric and labelKeys attributes
func (f *PrometheusCounterFactory) Create(metric string, labelKeys []string) CounterVec {
	p := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metric,
			Help: " ",
		},
		labelKeys,
	)
	prometheus.Register(p)

	return p
}

// NewPrometheus creates new prometheus incrementer instance
func NewPrometheus(counterFactory CounterFactory) *Prometheus {
	return &Prometheus{counter: nil, counterFactory: counterFactory}
}

// Increment increments metric in prometheus
func (i *Prometheus) Increment(metric string) {
	if i.counter == nil {
		i.counter = i.counterFactory.Create(metric, []string{})
	}

	i.counter.WithLabelValues().Inc()
}

// IncrementWithLabels increments metric in prometheus with defined labels
func (i *Prometheus) IncrementWithLabels(metric string, labels map[string]string) {
	var keys []string
	var values []string

	for key, value := range labels {
		keys = append(keys, key)
		values = append(values, value)
	}

	if i.counter == nil {
		i.counter = i.counterFactory.Create(metric, keys)
	}

	i.counter.WithLabelValues(values...).Inc()
}

// IncrementN increments metric by n in prometheus
func (i *Prometheus) IncrementN(metric string, n int) {
	if i.counter == nil {
		i.counter = i.counterFactory.Create(metric, []string{})
	}

	i.counter.WithLabelValues().Add(float64(n))
}

// IncrementNWithLabels increments metric by n in prometheus with defined labels
func (i *Prometheus) IncrementNWithLabels(metric string, n int, labels map[string]string) {
	keys := []string{}
	values := []string{}

	for key, value := range labels {
		keys = append(keys, key)
		values = append(values, value)
	}

	if i.counter == nil {
		i.counter = i.counterFactory.Create(metric, keys)
	}

	i.counter.WithLabelValues(values...).Add(float64(n))
}

// IncrementAll increments all metrics for given bucket in prometheus
func (i *Prometheus) IncrementAll(b bucket.Bucket) {

}

// IncrementAllN increments all metrics for given bucket in prometheus
func (i *Prometheus) IncrementAllN(b bucket.Bucket, n int) {
	incrementAllN(i, b, n)
}
