package incrementer

import (
	"github.com/hellofresh/stats-go/bucket"
	log "github.com/sirupsen/logrus"
)

// Log struct is Incrementer interface implementation that writes all metrics to log
type Log struct{}

// Increment writes given metric to log
func (i *Log) Increment(metric string) {
	log.WithField("metric", metric).Debug("Muted stats counter increment")
}

// IncrementN writes given metric to log
func (i *Log) IncrementN(metric string, n int) {
	log.WithField("metric", metric).WithField("n", n).Debug("Muted stats counter increment by n")
}

// IncrementAll writes all metrics for given bucket to log
func (i *Log) IncrementAll(b bucket.Bucket) {
	incrementAll(i, b)
}

// IncrementAllN writes all metrics for given bucket to log
func (i *Log) IncrementAllN(b bucket.Bucket, n int) {
	incrementAllN(i, b, n)
}
