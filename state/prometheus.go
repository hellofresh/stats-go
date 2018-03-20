package state

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus struct is State interface implementation that writes all states
type Prometheus struct {
	gauge        GaugeVec
	gaugeFactory GaugeFactory
}

// GaugeVec interface for gauge vectors in prometheus backend
type GaugeVec interface {
	GetMetricWithLabelValues(lvs ...string) (prometheus.Gauge, error)
	GetMetricWith(labels prometheus.Labels) (prometheus.Gauge, error)
	WithLabelValues(lvs ...string) prometheus.Gauge
	With(labels prometheus.Labels) prometheus.Gauge
}

// GaugeFactory interface for making new GaugeVec instances
type GaugeFactory interface {
	Create(metric string, labelKeys []string) GaugeVec
}

// PrometheusGaugeFactory implements GaugeFactory interface
type PrometheusGaugeFactory struct {
}

// NewPrometheusGaugeFactory returns new PrometheusGaugeFactory instance
func NewPrometheusGaugeFactory() *PrometheusGaugeFactory {
	return &PrometheusGaugeFactory{}
}

// Create method returns new GaugeVec instance with metric and labelKeys attributes
func (f *PrometheusGaugeFactory) Create(metric string, labelKeys []string) GaugeVec {
	p := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metric,
			Help: " ",
		},
		labelKeys,
	)
	prometheus.Register(p)

	return p
}

// NewPrometheus creates new prometheus state instance
func NewPrometheus(gaugeFactory GaugeFactory) *Prometheus {
	return &Prometheus{gauge: nil, gaugeFactory: gaugeFactory}
}

// Set sets metric state
func (s *Prometheus) Set(metric string, n int) {
	if s.gauge == nil {
		s.gauge = s.gaugeFactory.Create(metric, []string{})
	}
	s.gauge.WithLabelValues().Add(float64(n))
}

// SetWithLabels sets metric state
func (s *Prometheus) SetWithLabels(metric string, n int, labelNames, labelValues []string) {
	if s.gauge == nil {
		s.gauge = s.gaugeFactory.Create(metric, labelNames)
	}
	s.gauge.WithLabelValues(labelValues...).Add(float64(n))
}
