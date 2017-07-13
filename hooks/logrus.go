package hooks

import (
	"github.com/hellofresh/stats-go"
	"github.com/hellofresh/stats-go/bucket"
	"github.com/sirupsen/logrus"
)

// LogrusHook is logrus hook for gathering stats on logged error
type LogrusHook struct {
	statsClient stats.Client
	section     string
}

// NewLogrusHook creates a stats logger.
func NewLogrusHook(statsClient stats.Client, section string) *LogrusHook {
	return &LogrusHook{statsClient: statsClient, section: section}
}

// Levels is logrus.Hook method implementation
func (h *LogrusHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
}

// Fire is logrus.Hook method implementation
func (h *LogrusHook) Fire(e *logrus.Entry) error {
	h.statsClient.TrackMetric(h.section, bucket.MetricOperation{e.Level.String()})
	return nil
}
