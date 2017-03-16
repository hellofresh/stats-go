package stats

import (
	log "github.com/Sirupsen/logrus"
)

// LogIncrementer struct is Incrementer interface implementation that writes all metrics to log
type LogIncrementer struct{}

// Increment writes given metric to log
func (t *LogIncrementer) Increment(bucket string) {
	log.WithField("bucket", bucket).Debug("Muted stats counter increment")
}
