package timer

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus struct is Timer interface implementation that writes all timings
type Prometheus struct {
	timerStart    time.Time
	histogramVecs map[string]*prometheus.HistogramVec
}

// StartAt starts timer at a given time
func (t *Prometheus) StartAt(s time.Time) Timer {
	t.timerStart = s
	return t
}

// NewPrometheus creates new prometheus timer instance
func NewPrometheus() *Prometheus {
	return &Prometheus{histogramVecs: map[string]*prometheus.HistogramVec{}}
}

// Start starts timer
func (t *Prometheus) Start() Timer {
	t.timerStart = time.Now()
	return t
}

// Finish writes elapsed time for metric to prometheus
func (t *Prometheus) Finish(bucket string) {
	t.FinishWithLabels(bucket, map[string]string{})
}

// FinishWithLabels writes elapsed time for metric to prometheus
func (t *Prometheus) FinishWithLabels(bucket string, labels map[string]string) {
	keys := []string{}
	values := []string{}

	for key, value := range labels {
		keys = append(keys, key)
		values = append(values, value)
	}

	duration := time.Since(t.timerStart)

	if _, ok := t.histogramVecs[bucket]; !ok {
		t.histogramVecs[bucket] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: bucket + "_seconds",
			Help: " ",
		}, keys)
		prometheus.Register(t.histogramVecs[bucket])
	}
	t.histogramVecs[bucket].WithLabelValues(values...).Observe(duration.Seconds())
}
