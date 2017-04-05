package stats

import (
	log "github.com/Sirupsen/logrus"
)

// LogIncrementer struct is Incrementer interface implementation that writes all metrics to log
type LogIncrementer struct{}

// Increment writes given metric to log
func (i *LogIncrementer) Increment(metric string) {
	log.WithField("metric", metric).Debug("Muted stats counter increment")
}

// IncrementN writes given metric to log
func (i *LogIncrementer) IncrementN(metric string, n int) {
	log.WithField("metric", metric).WithField("n", n).Debug("Muted stats counter increment by n")
}

// IncrementAll writes all metrics for given bucket to log
func (i *LogIncrementer) IncrementAll(b Bucket) {
	incrementAll(i, b)
}

// IncrementAllN writes all metrics for given bucket to log
func (i *LogIncrementer) IncrementAllN(b Bucket, n int) {
	incrementAllN(i, b, n)
}
