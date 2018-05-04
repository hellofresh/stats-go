package incrementer

import (
	"sync"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus struct is Incrementer interface implementation that writes all metrics to Prometheus
type Prometheus struct {
	sync.Mutex

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

// PrometheusCounterFactory implements CounterFactory interface
type PrometheusCounterFactory struct{}

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

// PrometheusIncrementerFactory implements Factory interface
type PrometheusIncrementerFactory struct{}

// NewPrometheusIncrementerFactory returns new NewPrometheusIncrementerFactory instance
func NewPrometheusIncrementerFactory() *PrometheusIncrementerFactory {
	return &PrometheusIncrementerFactory{}
}

// Create method returns new Prometheus incrementer instance
func (p *PrometheusIncrementerFactory) Create() Incrementer {
	return NewPrometheus(NewPrometheusCounterFactory())
}

// NewPrometheus creates new prometheus incrementer instance
func NewPrometheus(counterFactory CounterFactory) *Prometheus {
	return &Prometheus{counter: nil, counterFactory: counterFactory}
}

// Increment increments metric in prometheus
func (i *Prometheus) Increment(metric string, labels ...map[string]string) {
	var labelNames []string
	var labelValues []string

	if labels != nil {
		for k, v := range labels[0] {
			labelNames = append(labelNames, k)
			labelValues = append(labelValues, v)
		}
	}

	i.Lock()
	defer i.Unlock()

	if i.counter == nil {
		i.counter = i.counterFactory.Create(metric, labelNames)
	}

	i.counter.WithLabelValues(labelValues...).Inc()
}

// IncrementN increments metric by n in prometheus
func (i *Prometheus) IncrementN(metric string, n int, labels ...map[string]string) {
	var labelNames []string
	var labelValues []string

	if labels != nil {
		for k, v := range labels[0] {
			labelNames = append(labelNames, k)
			labelValues = append(labelValues, v)
		}
	}

	i.Lock()
	defer i.Unlock()

	if i.counter == nil {
		i.counter = i.counterFactory.Create(metric, labelNames)
	}

	i.counter.WithLabelValues(labelValues...).Add(float64(n))
}

// IncrementAll increments all metrics for given bucket in prometheus
func (i *Prometheus) IncrementAll(b bucket.Bucket) {

}

// IncrementAllN increments all metrics for given bucket in prometheus
func (i *Prometheus) IncrementAllN(b bucket.Bucket, n int) {
	incrementAllN(i, b, n)
}
